package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type CallbackFun func(data interface{})

type PriorityQueue struct {
	highQueue   chan interface{}
	normalQueue chan interface{}
	callBack    CallbackFun
	jobWaiter   sync.WaitGroup
	count       int64
	cancel      context.CancelFunc
}

// NewPriority 创建一个优先级队列，params第一个是thread
func NewPriority(ctx context.Context, size int, fun CallbackFun, params ...int) *PriorityQueue {
	thread := 1
	if len(params) != 0 {
		thread = params[0]
	}
	q := &PriorityQueue{}
	q.init(ctx, size, thread, fun)
	return q
}

func (q *PriorityQueue) init(parenCtx context.Context, size, thread int, fun CallbackFun) {
	q.callBack = fun
	ctx, cancel := context.WithCancel(parenCtx)
	q.cancel = cancel
	q.highQueue = make(chan interface{}, size)
	q.normalQueue = make(chan interface{}, size)
	q.jobWaiter.Add(thread)
	for i := 0; i < thread; i++ {
		go q.run(ctx)
	}
}

func (q *PriorityQueue) Push(value interface{}) {
	atomic.AddInt64(&q.count, 1)
	q.normalQueue <- value
}
func (q *PriorityQueue) Front(value interface{}) {
	atomic.AddInt64(&q.count, 1)
	q.highQueue <- value
}

func (q *PriorityQueue) Count() int64 {
	return atomic.LoadInt64(&q.count)
}

func (q *PriorityQueue) Destroy() {
	if q.cancel != nil {
		q.cancel()
		q.cancel = nil
	}
	q.jobWaiter.Wait()
	close(q.highQueue)
	close(q.normalQueue)
}

func (q *PriorityQueue) run(ctx context.Context) {
	defer func() {
		q.jobWaiter.Done()
	}()
	for {
		select {
		case data := <-q.highQueue:
			atomic.AddInt64(&q.count, -1)
			q.callBack(data)
		case data := <-q.normalQueue:
			atomic.AddInt64(&q.count, -1)
			q.callBack(data)
		case <-time.After(time.Millisecond):
			break
		case <-ctx.Done():
			return
		}
	}
}
