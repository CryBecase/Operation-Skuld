package mcache

import (
	"fmt"
	"sync"
)

// Item 内存缓存单个实例
type Item interface{}

// MemCache 内存缓存结构
type MemCache struct {
	items map[string]Item
	mu    sync.RWMutex
}

// Set Add an item to the MemCache, replacing any existing item.
func (m *MemCache) Set(k string, x interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.set(k, x)
}

func (m *MemCache) set(k string, x interface{}) {
	m.items[k] = x
}

// Add an item to the MemCache only if an item doesn't already exist for the given
func (m *MemCache) Add(k string, x interface{}) error {
	m.mu.Lock()
	_, found := m.get(k)
	if found {
		m.mu.Unlock()
		return fmt.Errorf("Item %s already exists", k)
	}
	m.set(k, x)
	m.mu.Unlock()
	return nil
}

// Replace Set a new value for the MemCache key only if it already exists, and the existing
func (m *MemCache) Replace(k string, x interface{}) error {
	m.mu.Lock()
	_, found := m.get(k)
	if !found {
		m.mu.Unlock()
		return fmt.Errorf("Item %s doesn't exist", k)
	}
	m.set(k, x)
	m.mu.Unlock()
	return nil
}

// Get an item from the MemCache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (m *MemCache) Get(k string) (interface{}, bool) {
	m.mu.RLock()
	item, found := m.items[k]
	if !found {
		m.mu.RUnlock()
		return nil, false
	}
	m.mu.RUnlock()
	return item, true
}

func (m *MemCache) get(k string) (interface{}, bool) {
	item, found := m.items[k]
	if !found {
		return nil, false
	}
	return item, true
}

// Delete an item from the MemCache. Does nothing if the key is not in the MemCache.
func (m *MemCache) Delete(k string) {
	m.mu.Lock()
	m.delete(k)
	m.mu.Unlock()
}

func (m *MemCache) delete(k string) {
	delete(m.items, k)
}

// Items Copies all items in the MemCache into a new map and returns it.
func (m *MemCache) Items() map[string]Item {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mm := make(map[string]Item, len(m.items))
	for k, v := range m.items {
		mm[k] = v
	}
	return mm
}

// ItemCount Returns the number of items in the MemCache. This may include items that have
func (m *MemCache) ItemCount() int {
	m.mu.RLock()
	n := len(m.items)
	m.mu.RUnlock()
	return n
}

// Flush Delete all items from the MemCache.
func (m *MemCache) Flush() {
	m.mu.Lock()
	m.items = map[string]Item{}
	m.mu.Unlock()
}

// NewMCache New Go MemCache
func NewMCache() *MemCache {
	return &MemCache{items: make(map[string]Item)}
}
