package main

import (
	"lrucached/server"
	"time"
)

func main() {
	cached := server.NewCacheServer(100000, 10*time.Minute, 2*time.Minute)
	cached.Init()
	cached.Run(":8383")
}
