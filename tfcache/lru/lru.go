package lru

import "container/list"

/*
	LRU
	最近最少使用，相对于仅考虑时间因素的 FIFO 和仅考虑访问频率的 LFU，
	LRU 算法可以认为是相对平衡的一种淘汰算法。LRU 认为，如果数据最近被访问过，
	那么将来被访问的概率也会更高。LRU 算法的实现非常简单，维护一个队列，如果某条记录被访问了
	则移动到队尾，那么队首则是最近最少访问的数据，淘汰该条记录即可。
*/

// Cache is a LRU cache. It is not safe for concurrent access.
// 实现方法，使用map 映射key 与节点的关系，方便快速找到节点
// 找到相关节点元素后，在在底层的双向链表中去操作节点
// 该方法通过map 实现了 O(1) 的节点查找， 使用链表实现了O(1) 的节点添加删除移动
type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	nowBytes  int64                         // 当前已使用的内存
	ll        *list.List                    // 底层结构，双向链表
	cache     map[string]*list.Element      // 键是字符串，值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) // 某条记录被移除时的回调函数，可以为 ni
}

// 双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射
type entry struct {
	Key   string
	value Value
}

// Value use Len to count how many bytes it takes
// 实现通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
// 查找主要有 2 个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾。
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 获取节点
	if ele, ok := c.cache[key]; ok {
		// 该节点存在，将该节点移动到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 返回查找到的值
		return kv.value, true
	}
	return
}

// RemoveOldest remove the oldest item
// 删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	// 获取队首节点
	ele := c.ll.Back()
	if ele != nil {
		// 从链表中删除该节点
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 从字典中从c.cache 删除该节点的映射关系
		delete(c.cache, kv.Key)
		// 更新当前所使用的内存从。nowBytes
		c.nowBytes -= int64(len(kv.Key)) + int64(kv.value.Len())
		// 如果存在回调函数，则调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.Key, kv.value)
		}
	}
}

// Add 新增修改元素
// 如果键存在，则更新对应节点的值，并将该节点移到队尾。
// 不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
// 更新 c.nbytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
func (c *Cache) Add(key string, value Value) {
	// 判断元素是否存在
	if ele, ok := c.cache[key]; ok {
		// 存在将该元素节点移动到链表队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 增加当前内存使用
		c.nowBytes += int64(value.Len()) - int64(kv.value.Len())
		// 更新节点原始值
		kv.value = value
	} else {
		// 元素不存在，新建节点，插入链表
		ele := c.ll.PushFront(&entry{key, value})
		// map 添加key与ele节点的映射关系
		c.cache[key] = ele
		// 计算当前内存使用
		c.nowBytes += int64(len(key)) + int64(value.Len())
	}
	// 如果超过最大内存使用率，进行淘汰节点
	for c.maxBytes != 0 && c.maxBytes < c.nowBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
// 实现 Len() 用来获取添加了多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
