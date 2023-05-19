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

package client

import (
	context "context"
	json "encoding/json"
	errors "errors"
	fmt "fmt"

	v1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	versioned "github.com/nubank/hringhorni/pkg/client/clientset/versioned"
	typedhringhorniv1alpha1 "github.com/nubank/hringhorni/pkg/client/clientset/versioned/typed/rollouts/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	discovery "k8s.io/client-go/discovery"
	dynamic "k8s.io/client-go/dynamic"
	rest "k8s.io/client-go/rest"
	injection "knative.dev/pkg/injection"
	dynamicclient "knative.dev/pkg/injection/clients/dynamicclient"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterClient(withClientFromConfig)
	injection.Default.RegisterClientFetcher(func(ctx context.Context) interface{} {
		return Get(ctx)
	})
	injection.Dynamic.RegisterDynamicClient(withClientFromDynamic)
}

// Key is used as the key for associating information with a context.Context.
type Key struct{}

func withClientFromConfig(ctx context.Context, cfg *rest.Config) context.Context {
	return context.WithValue(ctx, Key{}, versioned.NewForConfigOrDie(cfg))
}

func withClientFromDynamic(ctx context.Context) context.Context {
	return context.WithValue(ctx, Key{}, &wrapClient{dyn: dynamicclient.Get(ctx)})
}

// Get extracts the versioned.Interface client from the context.
func Get(ctx context.Context) versioned.Interface {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		if injection.GetConfig(ctx) == nil {
			logging.FromContext(ctx).Panic(
				"Unable to fetch github.com/nubank/hringhorni/pkg/client/clientset/versioned.Interface from context. This context is not the application context (which is typically given to constructors via sharedmain).")
		} else {
			logging.FromContext(ctx).Panic(
				"Unable to fetch github.com/nubank/hringhorni/pkg/client/clientset/versioned.Interface from context.")
		}
	}
	return untyped.(versioned.Interface)
}

type wrapClient struct {
	dyn dynamic.Interface
}

var _ versioned.Interface = (*wrapClient)(nil)

func (w *wrapClient) Discovery() discovery.DiscoveryInterface {
	panic("Discovery called on dynamic client!")
}

func convert(from interface{}, to runtime.Object) error {
	bs, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("Marshal() = %w", err)
	}
	if err := json.Unmarshal(bs, to); err != nil {
		return fmt.Errorf("Unmarshal() = %w", err)
	}
	return nil
}

// HringhorniV1alpha1 retrieves the HringhorniV1alpha1Client
func (w *wrapClient) HringhorniV1alpha1() typedhringhorniv1alpha1.HringhorniV1alpha1Interface {
	return &wrapHringhorniV1alpha1{
		dyn: w.dyn,
	}
}

type wrapHringhorniV1alpha1 struct {
	dyn dynamic.Interface
}

func (w *wrapHringhorniV1alpha1) RESTClient() rest.Interface {
	panic("RESTClient called on dynamic client!")
}

func (w *wrapHringhorniV1alpha1) Analysises(namespace string) typedhringhorniv1alpha1.AnalysisInterface {
	return &wrapHringhorniV1alpha1AnalysisImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "hringhorni.nu.dev",
			Version:  "v1alpha1",
			Resource: "analysises",
		}),

		namespace: namespace,
	}
}

type wrapHringhorniV1alpha1AnalysisImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typedhringhorniv1alpha1.AnalysisInterface = (*wrapHringhorniV1alpha1AnalysisImpl)(nil)

func (w *wrapHringhorniV1alpha1AnalysisImpl) Create(ctx context.Context, in *v1alpha1.Analysis, opts v1.CreateOptions) (*v1alpha1.Analysis, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "hringhorni.nu.dev",
		Version: "v1alpha1",
		Kind:    "Analysis",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Analysis{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Analysis, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Analysis{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AnalysisList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.AnalysisList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Analysis, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Analysis{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) Update(ctx context.Context, in *v1alpha1.Analysis, opts v1.UpdateOptions) (*v1alpha1.Analysis, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "hringhorni.nu.dev",
		Version: "v1alpha1",
		Kind:    "Analysis",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Analysis{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) UpdateStatus(ctx context.Context, in *v1alpha1.Analysis, opts v1.UpdateOptions) (*v1alpha1.Analysis, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "hringhorni.nu.dev",
		Version: "v1alpha1",
		Kind:    "Analysis",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Analysis{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1AnalysisImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}

func (w *wrapHringhorniV1alpha1) Rollouts(namespace string) typedhringhorniv1alpha1.RolloutInterface {
	return &wrapHringhorniV1alpha1RolloutImpl{
		dyn: w.dyn.Resource(schema.GroupVersionResource{
			Group:    "hringhorni.nu.dev",
			Version:  "v1alpha1",
			Resource: "rollouts",
		}),

		namespace: namespace,
	}
}

type wrapHringhorniV1alpha1RolloutImpl struct {
	dyn dynamic.NamespaceableResourceInterface

	namespace string
}

var _ typedhringhorniv1alpha1.RolloutInterface = (*wrapHringhorniV1alpha1RolloutImpl)(nil)

func (w *wrapHringhorniV1alpha1RolloutImpl) Create(ctx context.Context, in *v1alpha1.Rollout, opts v1.CreateOptions) (*v1alpha1.Rollout, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "hringhorni.nu.dev",
		Version: "v1alpha1",
		Kind:    "Rollout",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Create(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Rollout{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1RolloutImpl) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return w.dyn.Namespace(w.namespace).Delete(ctx, name, opts)
}

func (w *wrapHringhorniV1alpha1RolloutImpl) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	return w.dyn.Namespace(w.namespace).DeleteCollection(ctx, opts, listOpts)
}

func (w *wrapHringhorniV1alpha1RolloutImpl) Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Rollout, error) {
	uo, err := w.dyn.Namespace(w.namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Rollout{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1RolloutImpl) List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.RolloutList, error) {
	uo, err := w.dyn.Namespace(w.namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.RolloutList{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1RolloutImpl) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Rollout, err error) {
	uo, err := w.dyn.Namespace(w.namespace).Patch(ctx, name, pt, data, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Rollout{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1RolloutImpl) Update(ctx context.Context, in *v1alpha1.Rollout, opts v1.UpdateOptions) (*v1alpha1.Rollout, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "hringhorni.nu.dev",
		Version: "v1alpha1",
		Kind:    "Rollout",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).Update(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Rollout{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1RolloutImpl) UpdateStatus(ctx context.Context, in *v1alpha1.Rollout, opts v1.UpdateOptions) (*v1alpha1.Rollout, error) {
	in.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "hringhorni.nu.dev",
		Version: "v1alpha1",
		Kind:    "Rollout",
	})
	uo := &unstructured.Unstructured{}
	if err := convert(in, uo); err != nil {
		return nil, err
	}
	uo, err := w.dyn.Namespace(w.namespace).UpdateStatus(ctx, uo, opts)
	if err != nil {
		return nil, err
	}
	out := &v1alpha1.Rollout{}
	if err := convert(uo, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (w *wrapHringhorniV1alpha1RolloutImpl) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return nil, errors.New("NYI: Watch")
}
