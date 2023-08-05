package timer

import (
	"context"
	"time"
)

type Timer chan struct{}

func (t Timer) Close() {
	t <- struct{}{}
}

func StartTimeWithChan(callback func(), interval int, exitSignal Timer) {
	go func() {
		for {
			select {
			case <-time.After(time.Duration(interval) * time.Millisecond):
				callback()
			case <-exitSignal:
				return
			}
		}
	}()
}

func StartTimeWithContext(ctx context.Context, callback func(), interval int) {
	go func() {
		for {
			select {
			case <-time.After(time.Duration(interval) * time.Millisecond):
				callback()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func StartTime(callback func(), interval int) (exitSignal Timer) {
	exitSignal = make(Timer)
	go func() {
		for {
			select {
			case <-time.After(time.Duration(interval) * time.Millisecond):
				callback()
			case <-exitSignal:
				return
			}
		}
	}()
	return exitSignal
}
