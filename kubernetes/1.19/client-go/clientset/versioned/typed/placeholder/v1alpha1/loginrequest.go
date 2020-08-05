/*
Copyright 2020 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/suzerain-io/placeholder-name-api/pkg/apis/placeholder/v1alpha1"
	scheme "github.com/suzerain-io/placeholder-name/kubernetes/1.19/client-go/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// LoginRequestsGetter has a method to return a LoginRequestInterface.
// A group's client should implement this interface.
type LoginRequestsGetter interface {
	LoginRequests() LoginRequestInterface
}

// LoginRequestInterface has methods to work with LoginRequest resources.
type LoginRequestInterface interface {
	Create(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.CreateOptions) (*v1alpha1.LoginRequest, error)
	Update(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.UpdateOptions) (*v1alpha1.LoginRequest, error)
	UpdateStatus(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.UpdateOptions) (*v1alpha1.LoginRequest, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.LoginRequest, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.LoginRequestList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LoginRequest, err error)
	LoginRequestExpansion
}

// loginRequests implements LoginRequestInterface
type loginRequests struct {
	client rest.Interface
}

// newLoginRequests returns a LoginRequests
func newLoginRequests(c *PlaceholderV1alpha1Client) *loginRequests {
	return &loginRequests{
		client: c.RESTClient(),
	}
}

// Get takes name of the loginRequest, and returns the corresponding loginRequest object, and an error if there is any.
func (c *loginRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.LoginRequest, err error) {
	result = &v1alpha1.LoginRequest{}
	err = c.client.Get().
		Resource("loginrequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of LoginRequests that match those selectors.
func (c *loginRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.LoginRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.LoginRequestList{}
	err = c.client.Get().
		Resource("loginrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested loginRequests.
func (c *loginRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("loginrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a loginRequest and creates it.  Returns the server's representation of the loginRequest, and an error, if there is any.
func (c *loginRequests) Create(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.CreateOptions) (result *v1alpha1.LoginRequest, err error) {
	result = &v1alpha1.LoginRequest{}
	err = c.client.Post().
		Resource("loginrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(loginRequest).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a loginRequest and updates it. Returns the server's representation of the loginRequest, and an error, if there is any.
func (c *loginRequests) Update(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.UpdateOptions) (result *v1alpha1.LoginRequest, err error) {
	result = &v1alpha1.LoginRequest{}
	err = c.client.Put().
		Resource("loginrequests").
		Name(loginRequest.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(loginRequest).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *loginRequests) UpdateStatus(ctx context.Context, loginRequest *v1alpha1.LoginRequest, opts v1.UpdateOptions) (result *v1alpha1.LoginRequest, err error) {
	result = &v1alpha1.LoginRequest{}
	err = c.client.Put().
		Resource("loginrequests").
		Name(loginRequest.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(loginRequest).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the loginRequest and deletes it. Returns an error if one occurs.
func (c *loginRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("loginrequests").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *loginRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("loginrequests").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched loginRequest.
func (c *loginRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LoginRequest, err error) {
	result = &v1alpha1.LoginRequest{}
	err = c.client.Patch(pt).
		Resource("loginrequests").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
