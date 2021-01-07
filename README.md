# Lrucached

A basic lru cache service with an key-value associative array to store data.

## TODO List
- [x] Define internal data structure of cache
- [x] Define Item object
- [x] Implement basic methods Set, Get, Remove, Size
- [ ] Gracefully handle out of memory (OOM) exception when setting large amount of data
- [x] capacity and size attrs
- [ ] Set upper bound size for key
- [ ] impl mechanism to resize cap, might be overhead, but droping items until the cache have enough size to put new item, also isn't an option.
- [x] cache invalidation mechanism
- [x] Apply default lru eviction
- [ ] Item expiration
- [ ] Deletion expired cache items, impl clean up
- [ ] Cache Stats
- [ ] Logging actions and background jobs
- [ ] Benchmarking
- [x] Basic http-server implementation
- [ ] unload server side, set some of the work to client-side
- [ ] Implement item specific ttl independent from cache `defaultExpiration`
