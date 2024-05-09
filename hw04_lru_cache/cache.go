package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity  int
	queue     List
	items     map[Key]*ListItem
	itemsLock sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, newValue interface{}) bool {
	c.itemsLock.Lock()
	itemList, ok := c.items[key]

	if ok {
		itemList.Value = newValue
		c.queue.MoveToFront(itemList)
		c.itemsLock.Unlock()
		return true
	}

	if c.queue.Len() == c.capacity {
		lastItem := c.queue.Back()
		delete(c.items, lastItem.Key)
		c.queue.Remove(lastItem)
	}

	c.queue.PushFront(newValue)
	c.queue.Front().Key = key
	c.items[key] = c.queue.Front()
	c.itemsLock.Unlock()
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.itemsLock.Lock()
	itemList, ok := c.items[key]
	if ok {
		result := itemList.Value
		c.queue.MoveToFront(itemList)
		c.itemsLock.Unlock()
		return result, true
	}
	c.itemsLock.Unlock()
	return nil, false
}

func (c *lruCache) Clear() {
	clearList := NewList()
	clearItems := make(map[Key]*ListItem, c.capacity)
	c.itemsLock.Lock()
	c.queue, c.items = clearList, clearItems
	c.itemsLock.Unlock()
}
