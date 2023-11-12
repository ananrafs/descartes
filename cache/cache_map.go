package cache

type Cache map[string]bool

func (c *Cache) GetType() string {
	return "cache.map"
}

func NewCache() CacheItf {
	newCache := make(Cache, 0)
	return &newCache
}

func (c *Cache) SetCache(key string, value bool) {
	(*c)[key] = value
}

func (c *Cache) TryGet(key string, dest *bool) bool {
	if cache, ok := (*c)[key]; ok {
		*dest = cache
		return true
	}
	return false
}
