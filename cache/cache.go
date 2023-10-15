package cache

var (
	cacheInstanceMap map[string]CacheItf = make(map[string]CacheItf)
)

type CacheItf interface {
	SetCache(string, bool)
	TryGet(string, *bool) bool
	GetType() string
	New() CacheItf
}

func Init(cacheInstances ...CacheItf) {
	for _, cache := range cacheInstances {
		cacheInstanceMap[cache.GetType()] = cache
	}
}

func Get(cacheInstanceType string) CacheItf {

	cache, ok := cacheInstanceMap[cacheInstanceType]
	if ok {
		return cache.New()
	}

	return nil
}
