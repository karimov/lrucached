package cache

import (
	"container/list"
	"errors"
	"hash/fnv"
	"log"
	"sync"
	"time"
)

/*
Mainly memcached API like interfaces

1. Define cache datastructure which contains Items and item's keys
1.1. Define Item object
1.2. Set, Get, Remove ,ethods
1.2.1 Gracefully handle out of memory (OOM) exception when setting large amount of data
2. capacity and size attrs
2.1. Set upper bound size for key
3. impl mechanism to resize cap, might be overhead, but droping items until the cache have enough size to put new item, also isn't an option.
4. cache invalidation mechanism
5. Apply default lru eviction
6. Item expiration
6.1. Deletion expired cache items, impl clean up
7. Cache Stats
8. Logging actions and background jobs
9. Benchmarking
10. client-server architecture
10.1. unload server side, set some of the work to client-side

TODO:
1. Implement item specific ttl independent from cache `defaultExpiration`
*/

type Cache interface {
	// Sets new item to cache by provided key, or updates if such key's item cached already
	Set(key string, data []byte)
	// Gets item by provided key, returns nil if cache misses
	Get(key string) ([]byte, bool)
	// Removes the item from cache by provided key, does nothing if key not in the cache
	Remove(key string)
	// Return all cached items size in byte
	Size() uint64
}

type cache struct {
	sync.RWMutex                  //  syncranizing Get, Set, Remove methods
	cap          uint64           // Capacity of cache in byte
	size         uint64           // Size of all cached items in byte
	items        map[uint64]*Item // cached items, map key to tem
	hashList     *list.List       // list of hash keys, maintained from most recent to least recent used element
	timeToLive   time.Duration    // timeToLive expiration time of items
	cleanUpTime  time.Duration    // Cleanup window time after which expired items get removed
}

type Item struct {
	data       []byte        // raw data
	timeToLive time.Duration // expiration time of this item
	element    *list.Element // reference to item's key in the keyList
	size       uint64        // raw data size in byte
}

var (
	// ErrMaxCacheSize -
	ErrMaxCacheSize = errors.New("Maximum cache size reached")
)

// Size method returns raw data size in byte
func (it *Item) Size() uint64 {
	return uint64(len(it.data))
}

func (c *cache) Size() uint64 {
	c.Lock()
	defer c.Unlock()

	return c.size

}

// Set method stores new item to the cache, overwrites existing item
// item timeToLiv set as cache default timeToLive parameter
func (c *cache) Set(key string, data []byte) {
	c.Lock()
	defer c.Unlock()

	c.set(key, data)
}

// Gets the item from the cache,
// returns item if found or nil
func (c *cache) Get(key string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()

	return c.get(key)

}
func (c *cache) Remove(key string) {
	c.Lock()
	defer c.Unlock()

	hash := gethashkey(key)
	c.remove(hash)

}

func (c *cache) set(key string, data []byte) {
	size := uint64(len(data))
	hash := gethashkey(key)

	// check if item
	if err := c.checkSize(size); err != nil {
		log.Println(err)
		return
	}
	// remove first if exists such item
	c.remove(hash)
	// free up the cache size for new item to set,
	// by lru eviction policy
	c.ensureCapacity(size)

	c.items[hash] = &Item{
		data:       data,
		timeToLive: c.timeToLive,
		element:    c.setHashelement(hash),
		size:       size,
	}
	c.size += size
}

// Give a new item to store to the cache,
// check if item's size won't overflow the cache capacity.
func (c *cache) checkSize(size uint64) error {
	if c.cap <= size {
		return ErrMaxCacheSize
	}
	return nil
}

func (c *cache) get(key string) ([]byte, bool) {
	hash := gethashkey(key)
	item, found := c.items[hash]
	if !found {
		return nil, false
	}
	c.hashList.MoveToFront(item.element)
	return item.data, true
}

// A function to record the given hashkey and mark it as last to be evicted
func (c *cache) setHashelement(hashkey uint64) *list.Element {
	if item, ok := c.items[hashkey]; ok {
		c.hashList.MoveToFront(item.element)
		return item.element
	}
	return c.hashList.PushFront(hashkey)
}

// Given the need to add some number of new bytes to the cache,
// evict items according to the eviction policy until there is room.
// The caller should hold the cache lock.
func (c *cache) ensureCapacity(toAdd uint64) {
	mustRemove := int64(c.size+toAdd) - int64(c.cap)
	for mustRemove > 0 {
		hash := c.hashList.Back().Value.(uint64)
		mustRemove -= int64(c.items[hash].Size())
		c.remove(hash)
	}
}

// Remove the item associated with the given hash key.
// The caller should hold the cache lock.
func (c *cache) remove(hash uint64) {
	if item, ok := c.items[hash]; ok {
		delete(c.items, hash)
		c.size -= item.Size()
		c.hashList.Remove(item.element)
	}
}

// used FNV hash function
func hashFunc(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}

// returns 64-bit hash
func gethashkey(k string) uint64 {
	return hashFunc([]byte(k))
}

// NewCache returns new object of lrucached
func NewCache(capacity uint64, ttl, ttc time.Duration) Cache {
	return &cache{
		cap:         capacity,
		items:       map[uint64]*Item{},
		hashList:    list.New(),
		timeToLive:  ttl,
		cleanUpTime: ttc,
	}
}
