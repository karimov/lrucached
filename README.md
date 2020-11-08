# Lrucached

A basic lru cache service with an key-value associative array to store data.

## TODO List
- [ ] Define internal data set/get items
- [ ] Benchmarks

- [ ] Define cache datastructure which contains Items and item's keys
- [ ] Define Item object
- [ ] Set, Get, Remove ,ethods
- [ ] Gracefully handle out of memory (OOM) exception when setting large amount of data
- [ ] capacity and size attrs
- [ ] Set upper bound size for key
- [ ] impl mechanism to resize cap, might be overhead, but droping items until the cache have enough size to put new item, also isn't an option.
- [ ] cache invalidation mechanism
- [ ] Apply default lru eviction
- [ ] Item expiration
- [ ] Deletion expired cache items, impl clean up
- [ ] Cache Stats
- [ ] Logging actions and background jobs
- [ ] Benchmarking
- [ ] client-server architecture
- [ ] unload server side, set some of the work to client-side

TODO:
- [ ] Implement item specific ttl independent from cache `defaultExpiration`
