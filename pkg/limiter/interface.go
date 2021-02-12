package limiter

// Release func.
type Release func()

// Limiter interface.
type Limiter interface {
	// Allow should be checks if Take call available
	// if return false, Take will lock
	Allow() bool
	// Take should be reserves a space for call
	// returns a Release func for releasing space
	Take() Release
}
