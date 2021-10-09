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
// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/nistal97/crd_controller/pkg/api/tess.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CiConfigLister helps list CiConfigs.
// All objects returned here must be treated as read-only.
type CiConfigLister interface {
	// List lists all CiConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.CiConfig, err error)
	// CiConfigs returns an object that can list and get CiConfigs.
	CiConfigs(namespace string) CiConfigNamespaceLister
	CiConfigListerExpansion
}

// ciConfigLister implements the CiConfigLister interface.
type ciConfigLister struct {
	indexer cache.Indexer
}

// NewCiConfigLister returns a new CiConfigLister.
func NewCiConfigLister(indexer cache.Indexer) CiConfigLister {
	return &ciConfigLister{indexer: indexer}
}

// List lists all CiConfigs in the indexer.
func (s *ciConfigLister) List(selector labels.Selector) (ret []*v1.CiConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CiConfig))
	})
	return ret, err
}

// CiConfigs returns an object that can list and get CiConfigs.
func (s *ciConfigLister) CiConfigs(namespace string) CiConfigNamespaceLister {
	return ciConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CiConfigNamespaceLister helps list and get CiConfigs.
// All objects returned here must be treated as read-only.
type CiConfigNamespaceLister interface {
	// List lists all CiConfigs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.CiConfig, err error)
	// Get retrieves the CiConfig from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.CiConfig, error)
	CiConfigNamespaceListerExpansion
}

// ciConfigNamespaceLister implements the CiConfigNamespaceLister
// interface.
type ciConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CiConfigs in the indexer for a given namespace.
func (s ciConfigNamespaceLister) List(selector labels.Selector) (ret []*v1.CiConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CiConfig))
	})
	return ret, err
}

// Get retrieves the CiConfig from the indexer for a given namespace and name.
func (s ciConfigNamespaceLister) Get(name string) (*v1.CiConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("ciconfig"), name)
	}
	return obj.(*v1.CiConfig), nil
}
