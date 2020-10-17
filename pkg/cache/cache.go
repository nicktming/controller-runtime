package cache

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"context"
	toolscache "k8s.io/client-go/tools/cache"
	"time"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/cache/internal"
)

type Cache interface {
	client.Reader

	Informers
}

type Informers interface {
	GetInformer(ctx context.Context, obj client.Object) (Informer, error)

	Start(ctx context.Context) error

	WaitForCacheSync(ctx context.Context) bool
}

type Informer interface {
	AddEventHandler(handler toolscache.ResourceEventHandler)

	AddEventHandlerWithResyncPeriod(handler toolscache.ResourceEventHandler, resyncPeriod time.Duration)

	AddIndexers(indexers toolscache.Indexers) error

	HasSynced() bool
}

type Options struct {
	// Scheme is the scheme to use for mapping objects to GroupVersionKinds
	Scheme 		*runtime.Scheme
	// Mapper is the RESTMapper to use for mapping GroupVersionKinds to Resources
	Mapper 		meta.RESTMapper

	Resync 		*time.Duration

	Namespace 	string
}

var defaultResyncTime = 10 * time.Hour

func New(config *rest.Config, opts Options) (Cache, error) {
	opts, err := defaultOpts(config, opts)
	if err != nil {
		return nil, err
	}
	im := internal.NewInformersMap(config, opts.Scheme, opts.Mapper, *opts.Resync, opts.Namespace)
	return &informerCache{InformersMap: im}, nil
}

func defaultOpts(config *rest.Config, opts Options) (Options, error) {
	if opts.Scheme == nil {
		opts.Scheme = scheme.Scheme
	}

	if opts.Mapper == nil {
		var err error
		opts.Mapper, err = apiutil.NewDiscoveryRESTMapper(config)
		if err != nil {
			return opts, fmt.Errorf("could not create RESTMapper from config")
		}
	}

	if opts.Resync == nil {
		opts.Resync = &defaultResyncTime
	}

	return opts, nil
}