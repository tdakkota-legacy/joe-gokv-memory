package gokv

import "sync"

// Keys is abstraction for keys list
type Keys interface {
	OnAdd(key string)
	OnDelete(key string)
	Keys() []string
}

type mapKeys struct {
	keys map[string]struct{}
	m    sync.RWMutex
}

func (m mapKeys) OnAdd(key string) {
	m.m.Lock()
	defer m.m.Unlock()
	m.keys[key] = struct{}{}
}

func (m mapKeys) OnDelete(key string) {
	m.m.Lock()
	defer m.m.Unlock()
	delete(m.keys, key)
}

func (m mapKeys) Keys() []string {
	m.m.RLock()
	defer m.m.RUnlock()
	keys := make([]string, 0, len(m.keys))
	for k := range m.keys {
		keys = append(keys, k)
	}
	return keys
}
