package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int

	sync.Mutex
	queue List
	items map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.items[key]; ok {
		lItem := l.items[key]
		item := lItem.Value.(*cacheItem)
		item.value = value
		_ = l.queue.MoveToFront(lItem)
		return true
	}

	if l.queue.Len() == l.capacity {
		tail := l.queue.Back()
		delete(l.items, tail.Value.(*cacheItem).key)
		_ = l.queue.Remove(tail)
	}

	cItem := &cacheItem{
		key:   key,
		value: value,
	}
	newItem := l.queue.PushFront(cItem)
	l.items[key] = newItem
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.items[key]; ok {
		lItem := l.items[key]
		_ = l.queue.MoveToFront(lItem)
		return lItem.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
