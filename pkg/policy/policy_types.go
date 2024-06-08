package policy

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	multipleAttestorError = fmt.Errorf("mutliple attestor cannot be added in the same entry")
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:storageversion

// ImageVerificationPolicy defines rules to verify images used in matching resources
type ImageVerificationPolicy struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// ImageVerificationPolicy spec.
	Spec ImageVerificationPolicySpec `json:"spec" yaml:"spec"`
}

type ImageVerificationPolicySpec struct {
	Rules []ImageVerificationRule `json:"rules"`
}

type ImageVerificationRule struct {
	Name           string                `json:"name"`
	Match          any                   `json:"match"`
	ImageExtractor ImageExtractorConfigs `json:"imageExtractors"`
	Rules          VerificationRules     `json:"verify"`
}

type ImageExtractorConfigs map[string][]ImageExtractorConfig

type ImageExtractorConfig struct {
	// Path is the path to the object containing the image field in a custom resource.
	// It should be slash-separated. Each slash-separated key must be a valid YAML key or a wildcard '*'.
	// Wildcard keys are expanded in case of arrays or objects.
	Path string `json:"path" yaml:"path"`
	// Value is an optional name of the field within 'path' that points to the image URI.
	// This is useful when a custom 'key' is also defined.
	// +optional
	Value string `json:"value,omitempty" yaml:"value,omitempty"`
	// Name is the entry the image will be available under 'images.<name>' in the context.
	// If this field is not defined, image entries will appear under 'images.custom'.
	// +optional
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Key is an optional name of the field within 'path' that will be used to uniquely identify an image.
	// Note - this field MUST be unique.
	// +optional
	Key string `json:"key,omitempty" yaml:"key,omitempty"`
	// JMESPath is an optional JMESPath expression to apply to the image value.
	// This is useful when the extracted image begins with a prefix like 'docker://'.
	// The 'trim_prefix' function may be used to trim the prefix: trim_prefix(@, 'docker://').
	// Note - Image digest mutation may not be used when applying a JMESPAth to an image.
	// +optional
	JMESPath string `json:"jmesPath,omitempty" yaml:"jmesPath,omitempty"`
}

// VerificationRules is a set of VerificationPolicy
type VerificationRules []VerificationRule

// VerificationRule is a rule against which images are validated.
type VerificationRule struct {
	// ImageReferences is a list of matching image reference patterns. At least one pattern in the
	// list must match the image for the rule to apply. Each image reference consists of a registry
	// address, repository, image, and tag (defaults to latest). Wildcards ('*' and '?') are allowed.
	ImageReferences string `json:"imageReferences"`

	// Cosign is an array of attributes used to verify cosign signatures
	Cosign []*Cosign `json:"cosign,omitempty"`

	// Notary is an array of attributes used to verify notary signatures
	Notary []*Notary `json:"notary,omitempty"`
}

// Cosign is a set of attributes used to verify cosign signatures
type Cosign struct {
	Key                *Key           `json:"key,omitempty"`
	Keyless            *Keyless       `json:"keyless,omitempty"`
	Certificate        *Certificate   `json:"certificate,omitempty"`
	Rekor              *Rekor         `json:"rekor,omitempty"`
	CTLog              *CTLog         `json:"ctlog,omitempty"`
	SignatureAlgorithm string         `json:"signatureAlgorithm,omitempty"`
	Repository         string         `json:"repository,omitempty"`
	IgnoreTlog         bool           `json:"ignoreTlog"`
	IgnoreSCT          bool           `json:"ignoreSCT"`
	TSACertChain       string         `json:"tsaCertChain"`
	InToToAttestations []*Attestation `json:"intotoAttestations,omitempty"`
}

type Key struct {
	PublicKey string `json:"publicKey"`
}

type Keyless struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
	Root    string `json:"root"`
}

type Certificate struct {
	Cert      string `json:"cert"`
	CertChain string `json:"certChain"`
}

type Rekor struct {
	URL    string `json:"url"`
	PubKey string `json:"pubKey"`
}

type CTLog struct {
	PubKey string `json:"pubKey"`
}

// Notary is a set of attributes used to verify notary signatures
type Notary struct {
	Certs        string         `json:"certs"`
	Attestations []*Attestation `json:"attestations"`
}

type Attestation struct {
	Type string `json:"type"`
}

func (v *VerificationRule) Validate() error {
	for _, v := range v.Cosign {
		if v != nil {
			attestorAlreadyExists := false
			if v.Key != nil {
				if attestorAlreadyExists {
					return multipleAttestorError
				}
				attestorAlreadyExists = true
			}
			if v.Keyless != nil {
				if attestorAlreadyExists {
					return multipleAttestorError
				}
				attestorAlreadyExists = true
			}
			if v.Certificate != nil {
				if attestorAlreadyExists {
					return multipleAttestorError
				}
				attestorAlreadyExists = true
			}
		}
	}
	return nil
}
