package lrucache

import (
	"fmt"
	"time"
	"unsafe"
)

type LruCache struct {
	freeQueue *lruCacheQueue
	cache     *lruCacheQueue
	hasht     map[interface{}]interface{}
	key2node  map[interface{}]*lruCacheQueue
	node2key  map[string]interface{}
}

func New(size int64) *LruCache {
	if size < 1 {
		return nil
	}
	cache := new(LruCache)

	cache.freeQueue = queueInit(size)
	cache.cache = queueInit(0)
	cache.hasht = make(map[interface{}]interface{})
	cache.key2node = make(map[interface{}]*lruCacheQueue)
	cache.node2key = make(map[string]interface{})

	return cache
}

func ptr2str(ptr *lruCacheQueue) string {
	if ptr == nil {
		return ""
	}
	return fmt.Sprintf("%d", unsafe.Pointer(ptr))
}

func (c *LruCache) Get(key interface{}) (interface{}, bool) {
	var val interface{}
	var ok bool
	var node *lruCacheQueue

	val, ok = c.hasht[key]
	if !ok {
		return nil, false
	}
	node = c.key2node[key]
	if node == nil {
		return nil, false
	}
	queueRemove(node)
	queueInsertHead(c.cache, node)

	if (!node.expire.IsZero()) && node.expire.Sub(time.Now()) < 0 {
		return val, true
	}

	return val, false
}

func (c *LruCache) Set(key, value interface{}, ttl time.Duration) {
	var oldKey interface{}
	var ok bool
	var node *lruCacheQueue

	c.hasht[key] = value

	node, ok = c.key2node[key]
	if !ok {
		if queueIsEmpty(c.freeQueue) {
			node = queueLast(c.cache)
			oldKey = c.node2key[ptr2str(node)]
			delete(c.hasht, oldKey)
			delete(c.key2node, oldKey)

		} else {
			node = queueHead(c.freeQueue)
		}
		c.key2node[key] = node
		c.node2key[ptr2str(node)] = key
	}

	queueRemove(node)
	queueInsertHead(c.cache, node)

	if ttl >= 0 {
		node.expire = time.Now().Add(ttl)
	}
}

func (c *LruCache) Delete(key interface{}) {
	var node *lruCacheQueue

	delete(c.hasht, key)

	node = c.key2node[key]
	if node == nil {
		return
	}
	delete(c.node2key, ptr2str(node))
	delete(c.key2node, key)
	queueRemove(node)
	queueInsertTail(c.freeQueue, node)
}
