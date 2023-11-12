package cache

type NopCache struct{}

func (c *NopCache) GetType() string {
	return ""
}

func NewNopCache() CacheItf {
	return new(NopCache)
}

func (c *NopCache) SetCache(key string, value bool) {}

func (c *NopCache) TryGet(key string, dest *bool) bool {
	return false
}
