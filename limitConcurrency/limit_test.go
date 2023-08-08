package limitConcurrency

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	limit := NewLimit(2, 20)
	var id int32
	var job sync.WaitGroup
	job.Add(10 * 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			for j := 0; j < 10; j++ {
				limit.Request(func() {
					time.Sleep(time.Second)
					v := atomic.AddInt32(&id, 1)
					fmt.Println(index, v)
					job.Done()
				})
			}
		}(i)
	}
	job.Wait()
}
