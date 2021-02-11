package limiter

// concurrent implements of Limiter interface
// limits the count of calls.
type concurrent struct {
	limit int
	stack chan struct{}
}

// New create a new Limiter instance.
func New(limit int) Limiter {
	return &concurrent{
		limit: limit,
		stack: make(chan struct{}, limit),
	}
}

// Allow returns true if stack takes place.
func (lim *concurrent) Allow() bool {
	return lim.limit > len(lim.stack)
}

// Take takes up space in stack and return Release func.
func (lim *concurrent) Take() Release {
	lim.stack <- struct{}{}
	return lim.release
}

// release releases stack space.
func (lim *concurrent) release() {
	<-lim.stack
}
