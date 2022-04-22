package lru

import (
	"sync"

	"github.com/bpowers/approx-lru/simplelru"
)

// Cache is a thread-safe fixed size LRU cache.
type Cache struct {
	lock sync.Mutex
	lru  simplelru.LRU
	_    [16]byte
}

// New creates an LRU of the given size.
func New(size int) (*Cache, error) {
	return NewWithEvict(size, nil)
}

// NewWithEvict constructs a fixed size cache with the given eviction
// callback.
func NewWithEvict(size int, onEvicted func(key string, value interface{})) (*Cache, error) {
	lru, err := simplelru.NewLRU(size, onEvicted)
	if err != nil {
		return nil, err
	}
	c := &Cache{
		lru: *lru,
	}
	return c, nil
}

// Purge is used to completely clear the cache.
func (c *Cache) Purge() {
	c.lock.Lock()
	c.lru.Purge()
	c.lock.Unlock()
}

// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *Cache) Add(key string, value interface{}) (evicted bool) {
	c.lock.Lock()
	evicted = c.lru.Add(key, value)
	c.lock.Unlock()
	return evicted
}

// Get looks up a key's value from the cache.
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.Get(key)
	c.lock.Unlock()
	return value, ok
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (c *Cache) Contains(key string) bool {
	c.lock.Lock()
	containKey := c.lru.Contains(key)
	c.lock.Unlock()
	return containKey
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *Cache) Peek(key string) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.Peek(key)
	c.lock.Unlock()
	return value, ok
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache) ContainsOrAdd(key string, value interface{}) (ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	}
	evicted = c.lru.Add(key, value)
	return false, evicted
}

// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *Cache) PeekOrAdd(key string, value interface{}) (previous interface{}, ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	previous, ok = c.lru.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = c.lru.Add(key, value)
	return previous, false, evicted
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key string) (present bool) {
	c.lock.Lock()
	present = c.lru.Remove(key)
	c.lock.Unlock()
	return
}

// Resize changes the cache size.
func (c *Cache) Resize(size int) (evicted int) {
	c.lock.Lock()
	evicted = c.lru.Resize(size)
	c.lock.Unlock()
	return evicted
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.lock.Lock()
	length := c.lru.Len()
	c.lock.Unlock()
	return length
}
