package cache

var (
	cacheInstanceMap map[string]Factory = make(map[string]Factory)
)

type Factory func() CacheItf

type CacheItf interface {
	SetCache(string, bool)
	TryGet(string, *bool) bool
	GetType() string
}

func Init(cacheInstances ...Factory) {
	for _, cache := range cacheInstances {
		cacheInstanceMap[cache().GetType()] = cache
	}
}

func Get(cacheInstanceType string) CacheItf {
	cache, ok := cacheInstanceMap[cacheInstanceType]
	if ok {
		return cache()
	}

	return nil
}
