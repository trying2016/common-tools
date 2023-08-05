package queue

import (
	"context"
	"sync"
	"sync/atomic"
)

type Queue struct {
	queue     chan interface{}
	cancel    context.CancelFunc
	callBack  func(data interface{})
	jobWaiter sync.WaitGroup
	count     int64
}

// NewNormal 创建一个队列，params第一个表示 thread 数量
func NewNormal(ctx context.Context, size int, fun func(data interface{}), params ...int) *Queue {
	q := &Queue{}
	q.init(ctx, size, fun, params...)
	return q
}

func (q *Queue) init(parenCtx context.Context, size int, fun func(data interface{}), params ...int) {
	ctx, cancel := context.WithCancel(parenCtx)
	q.cancel = cancel

	q.callBack = fun
	q.queue = make(chan interface{}, size)
	thread := 1
	if len(params) > 0 {
		thread = params[0]
	}
	q.jobWaiter.Add(thread)
	for i := 0; i < thread; i++ {
		go q.run(ctx)
	}
}

func (q *Queue) Push(value interface{}) {
	atomic.AddInt64(&q.count, 1)
	q.queue <- value
}
func (q *Queue) Count() int64 {
	return atomic.LoadInt64(&q.count)
}

func (q *Queue) Destroy() {
	if q.cancel == nil {
		return
	}
	q.cancel()
	q.jobWaiter.Wait()
	close(q.queue)
	q.cancel = nil
}

func (q *Queue) run(ctx context.Context) {
	for {
		select {
		case data := <-q.queue:
			atomic.AddInt64(&q.count, -1)
			q.callBack(data)
		case <-ctx.Done():
			q.jobWaiter.Done()
			return
		}
	}
}
