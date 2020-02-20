package gokv

import (
	"github.com/philippgille/gokv/encoding"
	"github.com/tdakkota/joe-gokv-memory/test"
	"go.uber.org/zap"
	"testing"
)

func createMemory() (*MemoryStore, *test.MockStore) {
	m := test.NewMockStore(encoding.JSON)
	v, _ := NewMemory(m)
	return v, m
}

const key = "key"
const invalidKey = ""
const value = "value"

func TestMemoryStore_Set(t *testing.T) {
	t.Run("valid key", func(t *testing.T) {
		memory, store := createMemory()

		err := memory.Set(key, []byte(value))
		if err != nil {
			t.Error(err)
			return
		}

		if data, ok := store.Map()[key]; ok {
			var v []byte

			err := store.Codec.Unmarshal(data, &v)
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
	})

	t.Run("invalid key", func(t *testing.T) {
		memory, _ := createMemory()

		err := memory.Set(invalidKey, []byte(value))
		if err == nil {
			t.Error("error expected")
			return
		}
	})
}

func TestMemoryStore_Get(t *testing.T) {
	t.Run("valid key", func(t *testing.T) {
		memory, store := createMemory()

		data, err := store.Codec.Marshal([]byte(value))
		if err != nil {
			t.Error(err)
			return
		}

		store.Map()[key] = data

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
	})

	t.Run("invalid key", func(t *testing.T) {
		memory, _ := createMemory()

		_, ok, err := memory.Get(invalidKey)
		if ok || err == nil {
			t.Error("error expected")
			return
		}
	})
}

func TestMemoryStore_Delete(t *testing.T) {
	t.Run("valid key", func(t *testing.T) {
		memory, store := createMemory()

		data, err := store.Codec.Marshal([]byte(value))
		if err != nil {
			t.Error(err)
			return
		}

		store.Map()[key] = data

		ok, err := memory.Delete(key)
		if err != nil {
			t.Error(err)
			return
		}

		if !ok {
			t.Error("expected true")
		}
	})

	t.Run("invalid key", func(t *testing.T) {
		memory, _ := createMemory()

		ok, err := memory.Delete(invalidKey)
		if err == nil {
			t.Error("error expected")
		}

		if ok {
			t.Error("expected false")
		}
	})
}

func TestMemoryStore_Keys(t *testing.T) {
	memory, _ := createMemory()

	err := memory.Set(key, []byte(value))
	if err != nil {
		t.Error(err)
		return
	}

	keys, err := memory.Keys()
	if err != nil {
		t.Error(err)
		return
	}

	if len(keys) < 1 || keys[0] != key {
		t.Error("expected one key:", key)
	}
}

func TestMemoryStore_Close(t *testing.T) {
	memory, store := createMemory()

	err := memory.Close()
	if err != nil {
		t.Error(err)
		return
	}

	if !store.Closed() {
		t.Error("expected that Store is closed")
	}
}

func TestNewMemory(t *testing.T) {
	logger := zap.L()

	m, err := NewMemory(test.NewMockStore(encoding.JSON), WithLogger(logger))
	if err != nil {
		t.Error(err)
		return
	}

	if m.logger != logger {
		t.Error("expected equal")
		return
	}

	err = m.Set(key, []byte(value))
	if err != nil {
		t.Error(err)
		return
	}
}

func TestMemory(t *testing.T) {
	m := Memory(test.NewMockStore(encoding.JSON))
	if m == nil {
		t.Error("expected *joe.Module instance, not nil")
		return
	}
}
