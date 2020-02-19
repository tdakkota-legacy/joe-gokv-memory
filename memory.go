package gokv

import (
	"github.com/go-joe/joe"
	"github.com/philippgille/gokv"
	"go.uber.org/zap"
)

type memory struct {
	logger *zap.Logger
	store  gokv.Store
	// store
	keys Keys
}

// Memory creates a joe.Module
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

// NewMemory is memory struct constructor
func NewMemory(store gokv.Store, options ...Option) (*memory, error) {
	m := &memory{store: store}

	for _, opt := range options {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}

	if m.keys == nil {
		m.keys = mapKeys{}
	}

	if m.logger == nil {
		m.logger = zap.NewNop()
	}

	return m, nil
}

func (m memory) Set(key string, value []byte) error {
	m.keys.OnAdd(key)
	return m.store.Set(key, value)
}

func (m memory) Get(key string) (value []byte, ok bool, err error) {
	ok, err = m.store.Get(key, &value)
	return
}

func (m memory) Delete(key string) (bool, error) {
	m.keys.OnDelete(key)
	err := m.store.Delete(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m memory) Keys() ([]string, error) {
	return m.keys.Keys(), nil
}

func (m memory) Close() error {
	return m.store.Close()
}
