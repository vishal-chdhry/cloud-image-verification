package imageverifier

import (
	"context"
	"time"

	"github.com/kyverno/kyverno/pkg/clients/dclient"
	"github.com/kyverno/kyverno/pkg/config"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	"github.com/nirmata/json-image-verification/pkg/apis/v1alpha1"
	"github.com/nirmata/json-image-verification/pkg/policy"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type engine struct {
	client dclient.Interface
}

type Request struct {
	Policies []*v1alpha1.ImageVerificationPolicy
	Resource interface{}
}

type Response struct {
	Resource        interface{}
	PolicyResponses []PolicyResponse
}

type PolicyResponse struct {
	Policy        v1alpha1.ImageVerificationPolicy
	RuleResponses []RuleResponse
}

type RuleResponse struct {
	Rule               v1alpha1.ImageVerificationRule
	VerificationResult VerificationResult
}

type VerificationResult struct {
	Image               string
	VerificationOutcome VerificationOutcome
	// Error is only populated for ERROR verification outcome
	Error                 error
	VerificationResponses []VerificationResponse
}

type VerificationResponse struct {
	VerificationRule v1alpha1.VerificationRule
	Failures         []error
}

type VerificationOutcome string

const (
	PASS  VerificationOutcome = "PASS"
	SKIP  VerificationOutcome = "SKIP"
	FAIL  VerificationOutcome = "FAIL"
	ERROR VerificationOutcome = "ERROR"
)

func NewEngine(ctx context.Context, kubeClient kubernetes.Interface, dynamicClient dynamic.Interface) (*engine, error) {
	client, err := dclient.NewClient(ctx, dynamicClient, kubeClient, 15*time.Second)
	if err != nil {
		return nil, err
	}
	return &engine{
		client: client,
	}, nil
}

func NewEngineFromDClient(client dclient.Interface) *engine {
	return &engine{
		client: client,
	}
}

func (e *engine) Apply(request Request) Response {
	response := Response{
		Resource:        request.Resource,
		PolicyResponses: make([]PolicyResponse, len(request.Policies)),
	}
	jp := jmespath.New(config.NewDefaultConfiguration(false))
	jsonContext := enginecontext.NewContext(jp)
	for i, pol := range request.Policies {
		policyResponse := PolicyResponse{
			Policy:        *pol,
			RuleResponses: make([]RuleResponse, len(pol.Spec.Rules)),
		}
		for j, r := range pol.Spec.Rules {
			jsonContext.Checkpoint()
			defer jsonContext.Restore()
			ruleResponse := RuleResponse{
				Rule: r,
			}
			err := addResourceToJsonContext(jsonContext, request.Resource)
			if err != nil {
				ruleResponse.VerificationResult = VerificationResult{
					VerificationOutcome: ERROR,
					Error:               err,
				}
				policyResponse.RuleResponses[j] = ruleResponse
				continue
			}

			errs, err := policy.Match(context.Background(), r.Match, request.Resource)
			if err != nil {
				ruleResponse.VerificationResult = VerificationResult{
					VerificationOutcome: ERROR,
					Error:               err,
				}
				policyResponse.RuleResponses[j] = ruleResponse
				continue
			}
			if len(errs) > 0 {
				ruleResponse.VerificationResult = VerificationResult{
					VerificationOutcome: SKIP,
				}
				policyResponse.RuleResponses[j] = ruleResponse
				continue
			}

			images, err := policy.GetImages(request.Resource, r.ImageExtractor)
			if err != nil {
				ruleResponse.VerificationResult = VerificationResult{
					VerificationOutcome: ERROR,
					Error:               err,
				}
				policyResponse.RuleResponses[j] = ruleResponse
				continue
			}

			err = addImagesToJsonContext(jsonContext, images)
			if err != nil {
				ruleResponse.VerificationResult = VerificationResult{
					VerificationOutcome: ERROR,
					Error:               err,
				}
				policyResponse.RuleResponses[j] = ruleResponse
				continue
			}

			err = addContextEntriesToJsonContext(jsonContext, e.client, jp, r.Context)
			if err != nil {
				ruleResponse.VerificationResult = VerificationResult{
					VerificationOutcome: ERROR,
					Error:               err,
				}
				policyResponse.RuleResponses[j] = ruleResponse
				continue
			}

			rule, err := substituteVariablesInRule(r.DeepCopy(), jsonContext)
			if err != nil {
				ruleResponse.VerificationResult = VerificationResult{
					VerificationOutcome: ERROR,
					Error:               err,
				}
				policyResponse.RuleResponses[j] = ruleResponse
				continue
			}

			verifier := NewVerifier(rule.Rules, e.client, jsonContext, jp, rule.RequiredCount)
			for _, v := range images {
				result := verifier.Verify(v)
				ruleResponse.VerificationResult = result
			}
			policyResponse.RuleResponses[j] = ruleResponse
		}
		response.PolicyResponses[i] = policyResponse
	}
	return response
}
