package gokv

import (
	"github.com/go-joe/joe"
	"github.com/philippgille/gokv"
	"go.uber.org/zap"
	"sync"
)

type memory struct {
	logger *zap.Logger
	store  gokv.Store
	// store
	keys      map[string]struct{}
	keysMutex sync.RWMutex
}

// creates a joe.Module
func Memory(store gokv.Store) joe.Module {
	return joe.ModuleFunc(func(conf *joe.Config) error {
		mem, err := NewMemory(store, WithLogger(conf.Logger("gokv-memory")))
		if err != nil {
			return err
		}

		conf.SetMemory(mem)
		return nil
	})
}

// memory struct constructor
func NewMemory(store gokv.Store, options ...Option) (*memory, error) {
	m := &memory{
		keys:  map[string]struct{}{},
		store: store,
	}

	for _, opt := range options {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}

	if m.logger == nil {
		m.logger = zap.NewNop()
	}

	return m, nil
}

func (m memory) Set(key string, value []byte) error {
	m.keysMutex.Lock()
	defer m.keysMutex.Unlock()
	m.keys[key] = struct{}{}

	return m.store.Set(key, value)
}

func (m memory) Get(key string) (value []byte, ok bool, err error) {
	ok, err = m.store.Get(key, &value)
	return
}

func (m memory) Delete(key string) (bool, error) {
	m.keysMutex.Lock()
	defer m.keysMutex.Unlock()
	delete(m.keys, key)

	err := m.store.Delete(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m memory) Keys() ([]string, error) {
	m.keysMutex.RLock()
	defer m.keysMutex.RUnlock()
	keys := make([]string, 0, len(m.keys))
	for k := range m.keys {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m memory) Close() error {
	return m.store.Close()
}
