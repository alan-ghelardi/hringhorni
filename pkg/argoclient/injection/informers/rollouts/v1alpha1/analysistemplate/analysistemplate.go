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

// Code generated by injection-gen. DO NOT EDIT.

package analysistemplate

import (
	context "context"

	apisrolloutsv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	versioned "github.com/nubank/hringhorni/pkg/argoclient/clientset/versioned"
	v1alpha1 "github.com/nubank/hringhorni/pkg/argoclient/informers/externalversions/rollouts/v1alpha1"
	client "github.com/nubank/hringhorni/pkg/argoclient/injection/client"
	factory "github.com/nubank/hringhorni/pkg/argoclient/injection/informers/factory"
	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/argoclient/listers/rollouts/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Argoproj().V1alpha1().AnalysisTemplates()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx), resourceVersion: injection.GetResourceVersion(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1alpha1.AnalysisTemplateInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch github.com/nubank/hringhorni/pkg/argoclient/informers/externalversions/rollouts/v1alpha1.AnalysisTemplateInformer from context.")
	}
	return untyped.(v1alpha1.AnalysisTemplateInformer)
}

type wrapper struct {
	client versioned.Interface

	namespace string

	resourceVersion string
}

var _ v1alpha1.AnalysisTemplateInformer = (*wrapper)(nil)
var _ rolloutsv1alpha1.AnalysisTemplateLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apisrolloutsv1alpha1.AnalysisTemplate{}, 0, nil)
}

func (w *wrapper) Lister() rolloutsv1alpha1.AnalysisTemplateLister {
	return w
}

func (w *wrapper) AnalysisTemplates(namespace string) rolloutsv1alpha1.AnalysisTemplateNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace, resourceVersion: w.resourceVersion}
}

// SetResourceVersion allows consumers to adjust the minimum resourceVersion
// used by the underlying client.  It is not accessible via the standard
// lister interface, but can be accessed through a user-defined interface and
// an implementation check e.g. rvs, ok := foo.(ResourceVersionSetter)
func (w *wrapper) SetResourceVersion(resourceVersion string) {
	w.resourceVersion = resourceVersion
}

func (w *wrapper) List(selector labels.Selector) (ret []*apisrolloutsv1alpha1.AnalysisTemplate, err error) {
	lo, err := w.client.ArgoprojV1alpha1().AnalysisTemplates(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector:   selector.String(),
		ResourceVersion: w.resourceVersion,
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apisrolloutsv1alpha1.AnalysisTemplate, error) {
	return w.client.ArgoprojV1alpha1().AnalysisTemplates(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		ResourceVersion: w.resourceVersion,
	})
}
