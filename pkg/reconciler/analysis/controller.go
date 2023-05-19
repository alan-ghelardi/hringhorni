package analysis

import (
	"context"

	"knative.dev/pkg/configmap"

	"github.com/nubank/hringhorni/pkg/apis/config"
	analysisinformer "github.com/nubank/hringhorni/pkg/client/injection/informers/rollouts/v1alpha1/analysis"
	rolloutinformer "github.com/nubank/hringhorni/pkg/client/injection/informers/rollouts/v1alpha1/rollout"
	analysisreconciler "github.com/nubank/hringhorni/pkg/client/injection/reconciler/rollouts/v1alpha1/analysis"
	"github.com/nubank/hringhorni/pkg/kayenta"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(ctx context.Context, watcher configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	configStore := config.NewStore(logger.Named("configs"))
	configStore.WatchConfigs(watcher)

	analysisInformer := analysisinformer.Get(ctx)
	rolloutInformer := rolloutinformer.Get(ctx)

	reconciler := &Reconciler{
		rolloutLister: rolloutInformer.Lister(),
		analyzers:     []Analyzer{kayenta.Get(ctx)},
	}

	impl := analysisreconciler.NewImpl(ctx, reconciler, func(*controller.Impl) controller.Options {
		return controller.Options{ConfigStore: configStore}
	})

	logger.Info("Setting up event handlers")

	analysisInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
