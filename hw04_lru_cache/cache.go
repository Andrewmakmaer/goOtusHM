package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity      int
	queue         List
	items         map[Key]*ListItem
	valToKeyItems map[*ListItem]Key
	itemsLock     sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:      capacity,
		queue:         NewList(),
		items:         make(map[Key]*ListItem, capacity),
		valToKeyItems: make(map[*ListItem]Key, capacity),
	}
}

func (c *lruCache) Set(key Key, newValue interface{}) bool {
	c.itemsLock.RLock()
	itemList, ok := c.items[key]
	c.itemsLock.RUnlock()
	if ok {
		itemList.Value = newValue
		c.queue.MoveToFront(itemList)
		return true
	}
	if c.queue.Len() == c.capacity {
		lastItem := c.queue.Back()
		c.itemsLock.RLock()
		lastItemKey, ok := c.valToKeyItems[lastItem]
		c.itemsLock.RUnlock()
		if ok {
			c.itemsLock.Lock()
			delete(c.items, lastItemKey)
			c.queue.Remove(lastItem)
			delete(c.valToKeyItems, lastItem)
			c.itemsLock.Unlock()
		}
	}
	c.queue.PushFront(newValue)
	c.itemsLock.Lock()
	c.items[key] = c.queue.Front()
	c.valToKeyItems[c.queue.Front()] = key
	c.itemsLock.Unlock()
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.itemsLock.RLock()
	itemList, ok := c.items[key]
	c.itemsLock.RUnlock()
	if ok {
		result := itemList.Value
		c.queue.MoveToFront(itemList)
		return result, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	clearList := NewList()
	clearItems := make(map[Key]*ListItem, c.capacity+1)
	clearValToKeyItems := make(map[*ListItem]Key, c.capacity+1)
	c.itemsLock.Lock()
	c.queue, c.items, c.valToKeyItems = clearList, clearItems, clearValToKeyItems
	c.itemsLock.Unlock()
}
