package lru

import (
	"errors"
	"reflect"
	"sync"
)

const (
	defaultCapacity = 1024
)

var (
	mCache map[string]*lruCache
	// ErrKey no key in cache
	ErrKey error = errors.New("Not Found By Key")

	// ErrConvert the value can not convert to  the given type
	ErrConvert error = errors.New("Convert error")
)

// Cache ia a cache use lru to delete key when is full
type Cache struct {
	mu       sync.RWMutex
	capacity int
	count    int
	head     *lruCache
	tail     *lruCache
	node     *lruCache
	mCache   map[string]*lruCache
}

type lruCache struct {
	key   string
	value interface{}
	pre   *lruCache
	next  *lruCache
}

func newDefaultCache() (c *Cache) {
	return &Cache{
		capacity: defaultCapacity,
		mCache:   make(map[string]*lruCache),
	}
}

// NewCache return the *Cache with capacity n if n <= 0 return newDefaultCache
func NewCache(capacity int) (c *Cache) {
	if capacity <= 0 {
		return newDefaultCache()
	}
	return &Cache{
		capacity: capacity,
		mCache:   make(map[string]*lruCache),
	}
}

// FlushCache delete all data
func (c *Cache) FlushCache() {
	c.mu.Lock()
	c.count = 0
	c.node = nil
	c.head = nil
	c.tail = nil
	c.mCache = make(map[string]*lruCache)
	c.mu.Unlock()
}

// Set set key:value to cache ,if already exists move to head
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if node, ok := c.mCache[key]; ok {
		// move to head
		c.moveNodeToHead(node)
	} else {
		node := &lruCache{
			key:   key,
			value: value,
		}
		if c.count >= c.capacity {
			// delete tail node count = c.capacity
			c.deleteNode()
		} else {
			c.count++
		}
		// insert node to head
		c.insertNode(node)
		c.mCache[key] = node
	}
}

func (c *Cache) moveNodeToHead(node *lruCache) {
	tmp := c.tail.pre
	c.tail.next = c.head
	c.head.pre = c.tail
	c.head = c.tail
	c.tail = tmp
}

func (c *Cache) deleteNode() {
	tmp := c.tail
	pre := c.tail.pre
	pre.next = nil
	c.tail = pre
	delete(c.mCache, tmp.key)
}

func (c *Cache) insertNode(node *lruCache) {
	node.next = c.head
	if c.tail == nil {
		c.tail = node
	} else {
		c.head.pre = node
	}
	c.head = node
}

// GetBool get key's value and convert to bool, if failed or not exists return err
func (c *Cache) GetBool(key string) (v bool, err error) {
	value, err := c.Get(key)
	if err != nil {
		return
	}
	value2 := reflect.ValueOf(value)
	switch value2.Kind() {
	case reflect.Bool:
		v = value2.Bool()
	default:
		err = ErrConvert
	}
	return
}

// GetInt get value of int64,if not int return error
func (c *Cache) GetInt(key string) (v int, err error) {
	value, err := c.getInt64(key)
	if err != nil {
		return
	}
	v = int(value)
	return
}

// GetInt8 like GetInt but return int8(value)
func (c *Cache) GetInt8(key string) (v int8, err error) {
	value, err := c.getInt64(key)
	if err != nil {
		return
	}
	v = int8(value)
	return
}

// GetInt16 like GetInt but return int16(value)
func (c *Cache) GetInt16(key string) (v int16, err error) {
	value, err := c.getInt64(key)
	if err != nil {
		return
	}
	v = int16(value)
	return
}

// GetInt32 like GetInt but return int32(value)
func (c *Cache) GetInt32(key string) (v int32, err error) {
	value, err := c.getInt64(key)
	if err != nil {
		return
	}
	v = int32(value)
	return
}

// GetInt64 like GetInt but return int64(value)
func (c *Cache) GetInt64(key string) (v int64, err error) {
	value, err := c.getInt64(key)
	if err != nil {
		return
	}
	v = value
	return
}

func (c *Cache) getInt64(key string) (v int64, err error) {
	value, err := c.Get(key)
	if err != nil {
		return
	}
	value2 := reflect.ValueOf(value)
	switch value2.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = value2.Int()
	default:
		err = ErrConvert
	}
	return
}

// GetUint get value of uint,if not exists or not uint  return error
func (c *Cache) GetUint(key string) (v uint, err error) {
	value, err := c.getUint64(key)
	if err != nil {
		return
	}
	v = uint(value)
	return
}

// GetUint8 like GetUint but return uint8
func (c *Cache) GetUint8(key string) (v uint8, err error) {
	value, err := c.getUint64(key)
	if err != nil {
		return
	}
	v = uint8(value)
	return
}

// GetUint16 like GetUint but return uint16
func (c *Cache) GetUint16(key string) (v uint16, err error) {
	value, err := c.getUint64(key)
	if err != nil {
		return
	}
	v = uint16(value)
	return
}

// GetUint32 like GetUint but return uint32
func (c *Cache) GetUint32(key string) (v uint32, err error) {
	value, err := c.getUint64(key)
	if err != nil {
		return
	}
	v = uint32(value)
	return
}

// GetUint64 like GetUint but return uint64
func (c *Cache) GetUint64(key string) (v uint64, err error) {
	value, err := c.getUint64(key)
	if err != nil {
		return
	}
	v = value
	return
}

func (c *Cache) getUint64(key string) (v uint64, err error) {
	value, err := c.Get(key)
	if err != nil {
		return
	}
	value2 := reflect.ValueOf(value)
	switch value2.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v = value2.Uint()
	default:
		err = ErrConvert
	}
	return
}

// GetStr an str value from cache,if not exists or can not convert to string,return nil err
func (c *Cache) GetStr(key string) (v string, err error) {
	value, err := c.Get(key)
	if err != nil {
		return
	}
	value2 := reflect.ValueOf(value)
	switch value2.Kind() {
	case reflect.String:
		v = value2.String()
	default:
		err = ErrConvert
	}
	return
}

// Get an item from cache ,if not exists return nil
func (c *Cache) Get(key string) (v interface{}, err error) {
	if node, ok := c.mCache[key]; ok {
		v = node.value
		c.moveNodeToHead(node)
	} else {
		err = ErrKey
	}
	return
}
