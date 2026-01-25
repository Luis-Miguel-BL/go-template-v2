package util

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"
)

func DoneChan(done chan struct{}) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		done <- struct{}{}
	}
}

func AwaitDone(done chan struct{}, timeout time.Duration) error {
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for done channel")
	}
}
