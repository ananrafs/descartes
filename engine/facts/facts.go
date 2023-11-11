package facts

import (
	"github.com/ananrafs/descartes/cache"
)

type FactsItf interface {
	GetMap() map[string]interface{}
	GetCacheInstance() cache.CacheItf
	SetCacheInstance(cache.CacheItf)
}

type Facts struct {
	Fields map[string]interface{} `json:"param"`
	cache  cache.CacheItf
}

func (f *Facts) GetMap() map[string]interface{} {
	return f.Fields
}

func (f *Facts) GetCacheInstance() cache.CacheItf {
	return f.cache
}

func (f *Facts) SetCacheInstance(cache cache.CacheItf) {
	f.cache = cache
}
