package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.Lock()
	defer cache.Unlock()

	if item, found := cache.items[key]; found {
		cacheItem := (item.Value).(*cacheItem)
		cacheItem.value = value
		cache.queue.MoveToFront(item)
		return true
	}

	cache.queue.Len() == cache.capacity {
		backItem := cache.queue.Back()
		cache.queue.Remove(backItem)
		delete(cache.items, backItem.Value.(*cacheItem).key)
	}

	cache.items[key] = cache.queue.PushFront(&cacheItem{key, value})
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.Lock()
	defer cache.Unlock()

	if item, found := cache.items[key]; found {
		cache.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (cache *lruCache) Clear() {
	cache.Lock()
	defer cache.Unlock()

	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
