package cache

import (
	"sigs.k8s.io/controller-runtime/pkg/cache/internal"
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"log"
)

type ErrCacheNotStarted struct{}

func (*ErrCacheNotStarted) Error() string {
	return "the cache is not started, can not read objects"
}

type informerCache struct {
	*internal.InformersMap
}


func (ip *informerCache) Get(ctx context.Context, key client.ObjectKey, out client.Object) error {
	gvk, err := apiutil.GVKForObject(out, ip.Scheme)
	if err != nil {
		return err
	}
	log.Printf("gvk: %v\n", gvk.String())
	started, cache, err := ip.InformersMap.Get(ctx, gvk, out)
	if err != nil {
		return err
	}
	if !started {
		return &ErrCacheNotStarted{}
	}
	return cache.Reader.Get(ctx, key, out)
}

func (ip *informerCache) GetInformer(ctx context.Context, obj client.Object) (Informer, error) {
	gvk, err := apiutil.GVKForObject(out, ip.Scheme)
	if err != nil {
		return err
	}
	log.Printf("gvk: %v\n", gvk.String())
	_, cache, err := ip.InformersMap.Get(ctx, gvk, out)
	if err != nil {
		return nil, err
	}
	return cache.Informer, nil
}

