package cache

import (
	"container/list"
	"errors"
	"hash/fnv"
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

TODO:
1. Implement item specific ttl independent from cache `defaultExpiration`
*/

type Cache interface {
	// Sets new item to cache by provided key, or updates if such key's item cached already
	Set(key string, data []byte)
	// Gets item by provided key, returns nil if cache misses
	Get(key string) []byte
	// Removes item from cache by provided key
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
	// ErrMaxMemorySize -
	ErrMaxMemorySize = errors.New("Maximum memory size reached")
)

// Size method returns raw data size in byte
func (it *Item) Size() uint64 {
	return uint64(len(it.data))
}

func (c *cache) Size() uint64 {

}

// Set method stores new item to the cache, overwrites existing item
// item timeToLiv set as cache default timeToLive parameter
func (c *cache) Set(key string, data []byte) {
	c.Lock()
	defer c.Unlock()

	c.set(key, data)

}
func (c *cache) Get(key uint64) []Item {

}
func (c *cache) Remove(uint64) {

}

func (c *cache) set(key string, data []byte) error {
	size := uint64(len(data))
	hash := hashkey(key)
	c.items[hash] = &Item{
		data:       data,
		timeToLive: c.timeToLive,
		element:    c.setHashelement(hash),
		size:       size,
	}
	c.size += size

}

// Give a new item to store to the chache,
// check if item's size won't overflow the chache capacity.
func (c *cache) checkSize(size uint64) error {
	if c.cap <= size {
		return ErrMaxMemorySize
	}
	return nil
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
func hashkey(k string) uint64 {
	return hashFunc([]byte(k))
}
