// Copyright 2020 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "go.pinniped.dev/generated/1.18/apis/crdpinniped/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CredentialIssuerConfigLister helps list CredentialIssuerConfigs.
type CredentialIssuerConfigLister interface {
	// List lists all CredentialIssuerConfigs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.CredentialIssuerConfig, err error)
	// CredentialIssuerConfigs returns an object that can list and get CredentialIssuerConfigs.
	CredentialIssuerConfigs(namespace string) CredentialIssuerConfigNamespaceLister
	CredentialIssuerConfigListerExpansion
}

// credentialIssuerConfigLister implements the CredentialIssuerConfigLister interface.
type credentialIssuerConfigLister struct {
	indexer cache.Indexer
}

// NewCredentialIssuerConfigLister returns a new CredentialIssuerConfigLister.
func NewCredentialIssuerConfigLister(indexer cache.Indexer) CredentialIssuerConfigLister {
	return &credentialIssuerConfigLister{indexer: indexer}
}

// List lists all CredentialIssuerConfigs in the indexer.
func (s *credentialIssuerConfigLister) List(selector labels.Selector) (ret []*v1alpha1.CredentialIssuerConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CredentialIssuerConfig))
	})
	return ret, err
}

// CredentialIssuerConfigs returns an object that can list and get CredentialIssuerConfigs.
func (s *credentialIssuerConfigLister) CredentialIssuerConfigs(namespace string) CredentialIssuerConfigNamespaceLister {
	return credentialIssuerConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CredentialIssuerConfigNamespaceLister helps list and get CredentialIssuerConfigs.
type CredentialIssuerConfigNamespaceLister interface {
	// List lists all CredentialIssuerConfigs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.CredentialIssuerConfig, err error)
	// Get retrieves the CredentialIssuerConfig from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.CredentialIssuerConfig, error)
	CredentialIssuerConfigNamespaceListerExpansion
}

// credentialIssuerConfigNamespaceLister implements the CredentialIssuerConfigNamespaceLister
// interface.
type credentialIssuerConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CredentialIssuerConfigs in the indexer for a given namespace.
func (s credentialIssuerConfigNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.CredentialIssuerConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CredentialIssuerConfig))
	})
	return ret, err
}

// Get retrieves the CredentialIssuerConfig from the indexer for a given namespace and name.
func (s credentialIssuerConfigNamespaceLister) Get(name string) (*v1alpha1.CredentialIssuerConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("credentialissuerconfig"), name)
	}
	return obj.(*v1alpha1.CredentialIssuerConfig), nil
}
