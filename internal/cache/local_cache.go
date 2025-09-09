package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type LocalCache struct {
	localCache *cache.Cache
	url        string
	defaultTTL time.Duration
}

var localCacheUse *LocalCache

func InitLocalCache(url string, defaultTTL time.Duration) {
	localCacheUse = &LocalCache{
		localCache: cache.New(time.Hour, 5*time.Second),
		url:        url,
		defaultTTL: time.Minute,
	}
}

func NewLocalCacheUseCase() *LocalCache {
	return localCacheUse
}

func (c *LocalCache) set(key string, value interface{}, ttl time.Duration) {
	c.localCache.Set(key, value, ttl)
}

func (c *LocalCache) get(key string) (interface{}, bool) {
	get, b := c.localCache.Get(key)
	if b {
		return get, b
	}
	return nil, b
}
