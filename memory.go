package gokv

import (
	"github.com/go-joe/joe"
	"github.com/philippgille/gokv"
	"go.uber.org/zap"
)

type MemoryStore struct {
	logger *zap.Logger
	store  gokv.Store
	// store
	keys Keys
}

// Memory creates a joe.Module
func Memory(store gokv.Store) joe.Module {
	return joe.ModuleFunc(func(conf *joe.Config) error {
		mem, err := NewMemory(store, WithLogger(conf.Logger("gokv-MemoryStore")))
		if err != nil {
			return err
		}

		conf.SetMemory(mem)
		return nil
	})
}

// NewMemory is MemoryStore struct constructor
func NewMemory(store gokv.Store, options ...Option) (*MemoryStore, error) {
	m := &MemoryStore{store: store}

	for _, opt := range options {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}

	if m.keys == nil {
		m.keys = &mapKeys{keys: map[string]struct{}{}}
	}

	if m.logger == nil {
		m.logger = zap.NewNop()
	}

	return m, nil
}

func (m MemoryStore) Set(key string, value []byte) error {
	m.keys.OnAdd(key)
	return m.store.Set(key, value)
}

func (m MemoryStore) Get(key string) (value []byte, ok bool, err error) {
	ok, err = m.store.Get(key, &value)
	return
}

func (m MemoryStore) Delete(key string) (bool, error) {
	m.keys.OnDelete(key)
	err := m.store.Delete(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m MemoryStore) Keys() ([]string, error) {
	return m.keys.Keys(), nil
}

func (m MemoryStore) Close() error {
	return m.store.Close()
}
