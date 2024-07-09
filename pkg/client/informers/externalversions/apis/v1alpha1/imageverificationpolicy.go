/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	apisv1alpha1 "github.com/nirmata/json-image-verification/pkg/apis/v1alpha1"
	versioned "github.com/nirmata/json-image-verification/pkg/client/clientset/versioned"
	internalinterfaces "github.com/nirmata/json-image-verification/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/nirmata/json-image-verification/pkg/client/listers/apis/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ImageVerificationPolicyInformer provides access to a shared informer and lister for
// ImageVerificationPolicies.
type ImageVerificationPolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ImageVerificationPolicyLister
}

type imageVerificationPolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewImageVerificationPolicyInformer constructs a new informer for ImageVerificationPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewImageVerificationPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredImageVerificationPolicyInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredImageVerificationPolicyInformer constructs a new informer for ImageVerificationPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredImageVerificationPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NirmataV1alpha1().ImageVerificationPolicies(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NirmataV1alpha1().ImageVerificationPolicies(namespace).Watch(context.TODO(), options)
			},
		},
		&apisv1alpha1.ImageVerificationPolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *imageVerificationPolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredImageVerificationPolicyInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *imageVerificationPolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apisv1alpha1.ImageVerificationPolicy{}, f.defaultInformer)
}

func (f *imageVerificationPolicyInformer) Lister() v1alpha1.ImageVerificationPolicyLister {
	return v1alpha1.NewImageVerificationPolicyLister(f.Informer().GetIndexer())
}
