package gokv

import "go.uber.org/zap"

type Option func(*memory) error

func WithLogger(logger *zap.Logger) Option {
	return func(m *memory) error {
		m.logger = logger
		return nil
	}
}
