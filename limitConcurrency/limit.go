// Package limitProcess Limit
/**
 * @Author: trying
 * @Description: 限制并发数
 * @File:  limit.go
 * @Version: 1.0.0
 * @Date: 2023/8/8 10:36
 */

package limitConcurrency

import (
	"context"
	"github.com/trying2016/common-tools/queue"
)

type Wait struct {
	fn func()
	ch chan struct{}
}

type LimitConcurrency struct {
	queue *queue.Queue
}

func NewLimit(limitSize, bufSize int) *LimitConcurrency {
	limit := &LimitConcurrency{}
	limit.queue = queue.NewNormal(context.Background(), bufSize, limit.run, limitSize)
	return limit
}

func (l *LimitConcurrency) Request(fn func()) {
	ch := make(chan struct{})

	l.queue.Push(&Wait{
		ch: ch,
		fn: fn,
	})
	<-ch
}

func (l *LimitConcurrency) run(v interface{}) {
	wait := v.(*Wait)
	wait.fn()
	wait.ch <- struct{}{}
}
