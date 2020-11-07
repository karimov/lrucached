package main

import (
	"fmt"
	"hash/fnv"
	"runtime"
	"time"

	"github.com/allegro/bigcache"
)

func main() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Println(m.TotalAlloc)
	key := `9select * from table where id == 9;`
	n := fingerprint([]byte(key))
	fmt.Println(n)
}

func fingerprint(b []byte) uint64 {
	hash := fnv.New64a()
	hash.Write(b)
	return hash.Sum64()
}

// func main() {
// 	PrintMemUsage()

// 	var size uint64 = 2000
// 	var key string
// 	var data item
// 	c := cache.New(5*time.Minute, 10*time.Minute)
// 	_c := lru.New(size)
// 	_bigc, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

// 	c1 := func(items int) {
// 		var i int = 1
// 		fmt.Println("go-datastructure: ", _c.Size())
// 		t0 := time.Now()
// 		// for i <= items {
// 		// if c.Size() > size || i > 1000 {
// 		// 	log.Println("run out of size or items: ", c.Size(), i)
// 		// 	break
// 		// }

// 		// iter 1
// 		key = strconv.Itoa(rand.Int())
// 		//data = make([]byte, rand.Intn(1000))
// 		data = make([]byte, 999)
// 		_c.Put(key, data)
// 		println(key, len(data))
// 		i++

// 		// iter 2
// 		key = strconv.Itoa(rand.Int())
// 		//data = make([]byte, rand.Intn(1000))
// 		data = make([]byte, 1000)
// 		_c.Put(key, data)
// 		println(key, len(data))

// 		// iter 3
// 		key = strconv.Itoa(rand.Int())
// 		//data = make([]byte, rand.Intn(1000))
// 		data = make([]byte, 2001)
// 		_c.Put(key, data)
// 		println(key, len(data))

// 		// }
// 		fmt.Println(time.Now().Sub(t0))
// 		PrintMemUsage()
// 		fmt.Println("go-datastructure: ", _c.Size())
// 	}

// 	c2 := func(items int) {
// 		var i int
// 		fmt.Println("go-cache: ", c.ItemCount())
// 		t0 := time.Now()
// 		for i <= items {

// 			// if c.Size() > size || i > 1000 {
// 			// 	log.Println("run out of size or items: ", c.Size(), i)
// 			// 	break
// 			// }
// 			key = strconv.Itoa(rand.Int())
// 			data = make([]byte, 1, rand.Intn(10000000))
// 			println(key, cap(data))
// 			c.Set(key, data, cache.NoExpiration)
// 			i++
// 		}
// 		fmt.Println(time.Now().Sub(t0))
// 		// PrintMemUsage()
// 		fmt.Println("go-cache: ", c.ItemCount())
// 	}

// 	c3 := func(items int) {
// 		var i int
// 		fmt.Println("bigcache:: ", _bigc.Capacity(), _bigc.Len())
// 		t0 := time.Now()
// 		for i <= items {

// 			// if c.Size() > size || i > 1000 {
// 			// 	log.Println("run out of size or items: ", c.Size(), i)
// 			// 	break
// 			// }
// 			key = strconv.Itoa(rand.Int())
// 			data = make([]byte, 1, rand.Intn(10000000))
// 			_bigc.Set(key, data)
// 			i++
// 		}
// 		fmt.Println(time.Now().Sub(t0))
// 		// PrintMemUsage()
// 		fmt.Println("bigcache: ", _bigc.Capacity(), _bigc.Len())
// 	}

// 	c1(10)
// 	// c2(1000)
// 	// c3(1000)
// 	println(c2)
// 	println(c3)

// }

type item []byte

func (i item) Size() uint64 {
	return uint64(len(i))
}

func testBigcache(items int) time.Duration {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	t0 := time.Now()
	for i := 0; i <= items; i++ {
		cache.Set("my-unique-key", []byte("value"))
		// entry, _ := cache.Get("my-unique-key")
		// fmt.Println(string(entry))
	}
	return time.Now().Sub(t0)

}

func memStat() {
	// Below is an example of using our PrintMemUsage() function
	// Print our starting memory usage (should be around 0mb)
	PrintMemUsage()

	var overall [][]int
	for i := 0; i < 4; i++ {

		// Allocate memory using make() and append to overall (so it doesn't get
		// garbage collected). This is to create an ever increasing memory usage
		// which we can track. We're just using []int as an example.
		a := make([]int, 0, 999999)
		overall = append(overall, a)

		// Print our memory usage at each interval
		PrintMemUsage()
		time.Sleep(time.Second)
	}

	// Clear our memory and print usage, unless the GC has run 'Alloc' will remain the same
	overall = nil
	PrintMemUsage()

	// Force GC to clear up, should see a memory drop
	runtime.GC()
	PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
