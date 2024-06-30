package lru

import "container/list"

// LRU  cache
type LRUCache struct {
	maxBytes int64 //最大字节数
	nBytes   int64 //现字节数
	myList   *list.List
	cache    map[string]*list.Element
	// 钩子函数
	onEvicated func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// 新建一个LRU缓存结构
func New(maxBytes int64, onEvicated func(string, Value)) *LRUCache {
	return &LRUCache{
		maxBytes:   maxBytes,
		myList:     list.New(),
		cache:      make(map[string]*list.Element),
		onEvicated: onEvicated,
	}
}

// 实现value接口 返回结构的长度
func (c *LRUCache) Len() int {
	return c.myList.Len()
}

// 增加一个节点
func (c *LRUCache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.myList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
	} else {
		ele := c.myList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// 访问一个节点
func (c *LRUCache) get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.myList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除最近最少未使用的节点（最后一个）
func (c *LRUCache) RemoveOldest() {
	ele := c.myList.Back()
	if ele != nil {
		c.myList.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicated != nil {
			c.onEvicated(kv.key, kv.value)
		}
	}
}
