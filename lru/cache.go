package lru

type LRUCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type Node struct {
	key   string
	value string
	prev  *Node
	next  *Node
}

type lruCache struct {
	capacity int
	cache    map[string]*Node
	head     *Node
	tail     *Node
}

func NewLruCache(capacity int) LRUCache {
	return &lruCache{
		capacity: capacity,
		cache:    make(map[string]*Node),
	}
}

func (lc *lruCache) Put(key, value string) {
	if node, exists := lc.cache[key]; exists {
		node.value = value
		lc.moveToFront(node)
		return
	}

	newNode := &Node{
		key:   key,
		value: value,
	}

	if len(lc.cache) >= lc.capacity {
		lc.removeNode(lc.tail)
		delete(lc.cache, lc.tail.key)
	}

	lc.addToFront(newNode)
	lc.cache[key] = newNode
}

func (lc *lruCache) Get(key string) (string, bool) {
	if node, exists := lc.cache[key]; exists {
		lc.moveToFront(node)
		return node.value, true
	}
	return "", false
}

func (lc *lruCache) moveToFront(node *Node) {
	lc.removeNode(node)
	lc.addToFront(node)
}

func (lc *lruCache) removeNode(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		lc.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		lc.tail = node.prev
	}

	node.prev = nil
	node.next = nil
}

func (lc *lruCache) addToFront(node *Node) {
	node.prev = nil
	node.next = lc.head

	if lc.head != nil {
		lc.head.prev = node
	}

	lc.head = node

	if lc.tail == nil {
		lc.tail = node
	}
}
