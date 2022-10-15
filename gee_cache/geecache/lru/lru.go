package lru

import "container/list"

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// entry is element in ll
type entry struct {
	key   string
	value Value
}

// Cache is a LRU cache. It is not thread-safe
type Cache struct {
	maxBytes  int64 // Max bytes allowed, '0' mean unlimited bytes
	nBytes    int64 // Current use bytes
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) // Optional func, executed when an entry is purged
}

// New is the constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get looks up a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key) + value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len returns the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
