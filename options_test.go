package gokv

import (
	"go.uber.org/zap"
	"testing"
)

func TestWithKeys(t *testing.T) {
	memory, _ := createMemory()
	k := &mapKeys{}

	err := WithKeys(k)(memory)
	if err != nil {
		t.Error()
	}

	if memory.keys != k {
		t.Error("expected equal")
	}
}

func TestWithLogger(t *testing.T) {
	memory, _ := createMemory()
	logger := zap.L()

	err := WithLogger(logger)(memory)
	if err != nil {
		t.Error()
	}

	if memory.logger != logger {
		t.Error("expected equal")
	}
}
