/*
Copyright 2023 The hringhorni Authors

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RolloutLister helps list Rollouts.
// All objects returned here must be treated as read-only.
type RolloutLister interface {
	// List lists all Rollouts in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Rollout, err error)
	// Rollouts returns an object that can list and get Rollouts.
	Rollouts(namespace string) RolloutNamespaceLister
	RolloutListerExpansion
}

// rolloutLister implements the RolloutLister interface.
type rolloutLister struct {
	indexer cache.Indexer
}

// NewRolloutLister returns a new RolloutLister.
func NewRolloutLister(indexer cache.Indexer) RolloutLister {
	return &rolloutLister{indexer: indexer}
}

// List lists all Rollouts in the indexer.
func (s *rolloutLister) List(selector labels.Selector) (ret []*v1alpha1.Rollout, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Rollout))
	})
	return ret, err
}

// Rollouts returns an object that can list and get Rollouts.
func (s *rolloutLister) Rollouts(namespace string) RolloutNamespaceLister {
	return rolloutNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RolloutNamespaceLister helps list and get Rollouts.
// All objects returned here must be treated as read-only.
type RolloutNamespaceLister interface {
	// List lists all Rollouts in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Rollout, err error)
	// Get retrieves the Rollout from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Rollout, error)
	RolloutNamespaceListerExpansion
}

// rolloutNamespaceLister implements the RolloutNamespaceLister
// interface.
type rolloutNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Rollouts in the indexer for a given namespace.
func (s rolloutNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Rollout, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Rollout))
	})
	return ret, err
}

// Get retrieves the Rollout from the indexer for a given namespace and name.
func (s rolloutNamespaceLister) Get(name string) (*v1alpha1.Rollout, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("rollout"), name)
	}
	return obj.(*v1alpha1.Rollout), nil
}
