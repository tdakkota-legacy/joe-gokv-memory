package test

import (
	"errors"
	"github.com/philippgille/gokv/encoding"
)

type MockStore struct {
	Codec  encoding.Codec
	m      map[string][]byte
	closed bool
}

func (store *MockStore) Map() map[string][]byte {
	return store.m
}

func (store *MockStore) Closed() bool {
	return store.closed
}

func NewMockStore(codec encoding.Codec) *MockStore {
	return &MockStore{Codec: codec, closed: false, m: map[string][]byte{}}
}

func (store *MockStore) Set(k string, v interface{}) error {
	if k == "" {
		return errors.New("invalid key")
	}

	d, err := store.Codec.Marshal(v)
	if err != nil {
		return err
	}

	store.m[k] = d

	return nil
}

func (store *MockStore) Get(k string, v interface{}) (bool, error) {
	if k == "" {
		return false, errors.New("invalid key")
	}

	d, ok := store.m[k]
	if !ok {
		return false, nil
	}
	return true, store.Codec.Unmarshal(d, v)
}

func (store *MockStore) Delete(k string) error {
	if k == "" {
		return errors.New("invalid key")
	}

	delete(store.m, k)
	return nil
}

func (store *MockStore) Close() error {
	store.closed = true
	return nil
}
