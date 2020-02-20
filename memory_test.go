package gokv

import (
	"github.com/philippgille/gokv/encoding"
	"testing"
)

type MockStore struct {
	m     map[string][]byte
	codec encoding.Codec
}

func (store *MockStore) Set(k string, v interface{}) error {
	d, err := store.codec.Marshal(v)
	if err != nil {
		return err
	}
	store.m[k] = d
	return nil
}

func (store *MockStore) Get(k string, v interface{}) (bool, error) {
	d, ok := store.m[k]
	if !ok {
		return false, nil
	}
	return true, store.codec.Unmarshal(d, v)
}

func (store *MockStore) Delete(k string) error {
	delete(store.m, k)
	return nil
}

func (store *MockStore) Close() error {
	return nil
}

func createMemory() (*MemoryStore, *MockStore) {
	m := &MockStore{codec: encoding.JSON, m: map[string][]byte{}}
	v, _ := NewMemory(m)
	return v, m
}

const key = "key"
const value = "value"

func TestMemoryStore_Set(t *testing.T) {
	memory, store := createMemory()

	err := memory.Set(key, []byte(value))
	if err != nil {
		t.Error(err)
		return
	}

	if data, ok := store.m[key]; ok {
		var v []byte

		err := store.codec.Unmarshal(data, &v)
		if err != nil {
			t.Error(err)
			return
		}

		if string(v) != value {
			t.Errorf("expected %s, got %s", value, string(v))
		}
	} else {
		t.Error("key doesn't exists")
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	memory, store := createMemory()

	data, err := store.codec.Marshal([]byte(value))
	if err != nil {
		t.Error(err)
		return
	}

	store.m[key] = data

	ok, err := memory.Delete(key)
	if err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Error("expected true")
	}
}

func TestMemoryStore_Get(t *testing.T) {
	memory, store := createMemory()

	data, err := store.codec.Marshal([]byte(value))
	if err != nil {
		t.Error(err)
		return
	}

	store.m[key] = data

	v, ok, err := memory.Get(key)
	if err != nil {
		t.Error(err)
		return
	}

	if !ok {
		t.Error("expected true")
	}

	if string(v) != value {
		t.Errorf("expected %s, got %s", value, string(v))
	}
}
