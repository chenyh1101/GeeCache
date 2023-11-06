package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(int64(c.cacheBytes), nil)
	}
	c.lru.Add(key, value)
}
func (c *cache) get(key string) (b ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if ele, ok := c.lru.Get(key); ok {
		return ele.(ByteView), ok
	}
	return
}