package main

// Pool is a worker pool.
type Pool chan bool

// NewPool creates a pool with n available workers.
func NewPool(n int) Pool {
	p := make(Pool, n)

	// add workers
	for i := 0; i < n; i++ {
		p <- true
	}

	return p
}
