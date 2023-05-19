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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAnalysises implements AnalysisInterface
type FakeAnalysises struct {
	Fake *FakeHringhorniV1alpha1
	ns   string
}

var analysisesResource = schema.GroupVersionResource{Group: "hringhorni.nu.dev", Version: "v1alpha1", Resource: "analysises"}

var analysisesKind = schema.GroupVersionKind{Group: "hringhorni.nu.dev", Version: "v1alpha1", Kind: "Analysis"}

// Get takes name of the analysis, and returns the corresponding analysis object, and an error if there is any.
func (c *FakeAnalysises) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Analysis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(analysisesResource, c.ns, name), &v1alpha1.Analysis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Analysis), err
}

// List takes label and field selectors, and returns the list of Analysises that match those selectors.
func (c *FakeAnalysises) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AnalysisList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(analysisesResource, analysisesKind, c.ns, opts), &v1alpha1.AnalysisList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AnalysisList{ListMeta: obj.(*v1alpha1.AnalysisList).ListMeta}
	for _, item := range obj.(*v1alpha1.AnalysisList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested analysises.
func (c *FakeAnalysises) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(analysisesResource, c.ns, opts))

}

// Create takes the representation of a analysis and creates it.  Returns the server's representation of the analysis, and an error, if there is any.
func (c *FakeAnalysises) Create(ctx context.Context, analysis *v1alpha1.Analysis, opts v1.CreateOptions) (result *v1alpha1.Analysis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(analysisesResource, c.ns, analysis), &v1alpha1.Analysis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Analysis), err
}

// Update takes the representation of a analysis and updates it. Returns the server's representation of the analysis, and an error, if there is any.
func (c *FakeAnalysises) Update(ctx context.Context, analysis *v1alpha1.Analysis, opts v1.UpdateOptions) (result *v1alpha1.Analysis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(analysisesResource, c.ns, analysis), &v1alpha1.Analysis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Analysis), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAnalysises) UpdateStatus(ctx context.Context, analysis *v1alpha1.Analysis, opts v1.UpdateOptions) (*v1alpha1.Analysis, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(analysisesResource, "status", c.ns, analysis), &v1alpha1.Analysis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Analysis), err
}

// Delete takes name of the analysis and deletes it. Returns an error if one occurs.
func (c *FakeAnalysises) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(analysisesResource, c.ns, name, opts), &v1alpha1.Analysis{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAnalysises) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(analysisesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.AnalysisList{})
	return err
}

// Patch applies the patch and returns the patched analysis.
func (c *FakeAnalysises) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Analysis, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(analysisesResource, c.ns, name, pt, data, subresources...), &v1alpha1.Analysis{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Analysis), err
}
