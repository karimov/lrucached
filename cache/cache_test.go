package cache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
	1. Unittesting
	1.1. given-when-then-and/or-then style-like testing
	2. Benchmark
*/

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func value() []byte {
	return make([]byte, 100)
}

func TestNewCacheSetGetRemove(t *testing.T) {
	// given
	var capacity uint64
	ttl := 10 * time.Microsecond
	ttc := 2 * time.Microsecond
	capacity = 100000000
	c := NewCache(capacity, ttl, ttc).(*cache) // capacity of cache 100 MB, with random value of ttl and ttc
	// when
	c.Set("foo", []byte("bar"))
	c.Set("bar", []byte("foo"))
	// then
	assert.NotNil(t, c, "New lrucached object")
	assert.Equal(t, c.cap, capacity, "Capacity must be 100 MB")
	assert.Equal(t, c.Size(), uint64(6))
	assert.NotNil(t, c.hashList)
	assert.Equal(t, c.hashList.Len(), 2)

	// when
	value, found := c.Get("foo")
	// then
	assert.NotNil(t, value, "value shouldn't be nil")
	assert.Equal(t, string(value), "bar")
	assert.True(t, found, "Item must be found")

	// when
	hash := gethashkey("foo")
	item, found := c.items[hash]
	//then
	assert.True(t, found)
	assert.Equal(t, item.element.Value, hash)
	assert.Equal(t, item.Size(), uint64(3))
	assert.Equal(t, item.size, item.Size())
	assert.Equal(t, item.data, []byte("bar"))

	//when
	_, found = c.Get("bar")
	hash = gethashkey("bar")
	item, _ = c.items[hash]
	//then
	assert.True(t, found)
	assert.Equal(t, item.element, c.hashList.Front())

	//when update
	c.Set("foo", []byte("tar"))
	value, found = c.Get("foo")
	hash = gethashkey("foo")
	item, _ = c.items[hash]
	// then
	assert.True(t, found)
	assert.Equal(t, value, []byte("tar"))
	assert.Equal(t, item.element, c.hashList.Front())

	// when remove
	c.Remove("foo")
	value, found = c.Get("foo")
	// then
	assert.False(t, found)
	assert.Nil(t, value)
	assert.Equal(t, uint64(3), c.Size())
}
