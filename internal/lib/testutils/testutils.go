package testutils

import (
	"errors"
	"sync"
	"time"
)

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil // completed normally
	case <-time.After(timeout):
		return errors.New("Timeout exceeded") // timed out
	}
}
