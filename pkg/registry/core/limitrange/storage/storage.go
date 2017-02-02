/*
Copyright 2015 The Kubernetes Authors.

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

package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/genericapiserver/registry/generic"
	genericregistry "k8s.io/kubernetes/pkg/genericapiserver/registry/generic/registry"
	"k8s.io/kubernetes/pkg/registry/cachesize"
	"k8s.io/kubernetes/pkg/registry/core/limitrange"
)

type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against horizontal pod autoscalers.
func NewREST(optsGetter generic.RESTOptionsGetter) *REST {
	store := &genericregistry.Store{
		Copier:      api.Scheme,
		NewFunc:     func() runtime.Object { return &api.LimitRange{} },
		NewListFunc: func() runtime.Object { return &api.LimitRangeList{} },
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*api.LimitRange).Name, nil
		},
		PredicateFunc:     limitrange.MatchLimitRange,
		QualifiedResource: api.Resource("limitranges"),
		WatchCacheSize:    cachesize.GetWatchCacheSizeByResource("limitranges"),

		CreateStrategy: limitrange.Strategy,
		UpdateStrategy: limitrange.Strategy,
		DeleteStrategy: limitrange.Strategy,
		ExportStrategy: limitrange.Strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: limitrange.GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		panic(err) // TODO: Propagate error up
	}
	return &REST{store}
}