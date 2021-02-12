package pool

import "context"

// Option function signature.
type Option func(p *Pool)

// WithContext returns Option function
// that change pool context to user's context.
func WithContext(ctx context.Context) Option {
	return func(p *Pool) {
		p.ctx = ctx
	}
}

// WithBuffer returns Option function
// that init buffered workers channel in Pool instance.
func WithBuffer(bufferSize int) Option {
	return func(p *Pool) {
		p.workers = make(chan Worker, bufferSize)
	}
}
