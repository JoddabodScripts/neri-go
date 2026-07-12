package nerimity

import (
	"container/list"
	"sync"
)

// cache is a goroutine-safe, optionally size-limited key/value store. It mirrors
// the JavaScript SDK's Collection: when a limit is set it behaves as an LRU
// cache (least-recently-used entries are evicted first, and reads promote an
// entry to most-recently-used). A limit of 0 means unbounded.
type cache[V any] struct {
	mu    sync.RWMutex
	limit int
	ll    *list.List // front = most recently used
	items map[string]*list.Element
}

type cacheEntry[V any] struct {
	key string
	val V
}

func newCache[V any](limit int) *cache[V] {
	return &cache[V]{
		limit: limit,
		ll:    list.New(),
		items: make(map[string]*list.Element),
	}
}

// get returns the value for key and whether it was present. On a hit, and when
// a limit is set, the entry is promoted to most-recently-used.
func (c *cache[V]) get(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	if c.limit > 0 {
		c.ll.MoveToFront(el)
	}
	return el.Value.(*cacheEntry[V]).val, true
}

// set inserts or updates key. When the cache is over its limit the
// least-recently-used entry is evicted.
func (c *cache[V]) set(key string, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		el.Value.(*cacheEntry[V]).val = val
		if c.limit > 0 {
			c.ll.MoveToFront(el)
		}
		return
	}
	el := c.ll.PushFront(&cacheEntry[V]{key: key, val: val})
	c.items[key] = el
	if c.limit > 0 && c.ll.Len() > c.limit {
		c.evictOldest()
	}
}

func (c *cache[V]) evictOldest() {
	el := c.ll.Back()
	if el == nil {
		return
	}
	c.ll.Remove(el)
	delete(c.items, el.Value.(*cacheEntry[V]).key)
}

// delete removes key if present.
func (c *cache[V]) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		c.ll.Remove(el)
		delete(c.items, key)
	}
}

// has reports whether key is present without affecting recency.
func (c *cache[V]) has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.items[key]
	return ok
}

// len returns the number of cached entries.
func (c *cache[V]) len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// values returns a snapshot slice of all cached values, most-recently-used
// first. Safe to range over while the cache is mutated concurrently.
func (c *cache[V]) values() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]V, 0, c.ll.Len())
	for el := c.ll.Front(); el != nil; el = el.Next() {
		out = append(out, el.Value.(*cacheEntry[V]).val)
	}
	return out
}

// clear removes every entry.
func (c *cache[V]) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ll.Init()
	c.items = make(map[string]*list.Element)
}
