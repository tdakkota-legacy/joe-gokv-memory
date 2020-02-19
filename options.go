package gokv

import "go.uber.org/zap"

// option func
type Option func(*memory) error

// sets logger
func WithLogger(logger *zap.Logger) Option {
	return func(m *memory) error {
		m.logger = logger
		return nil
	}
}
