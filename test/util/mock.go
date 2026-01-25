package util

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

func DoneChan(done chan struct{}) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		close(done)
	}
}

func AwaitDone(done chan struct{}, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-done:
		return nil
	case <-timer.C:
		return fmt.Errorf("timeout waiting for done channel")
	}
}
