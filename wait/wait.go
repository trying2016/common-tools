package wait

import (
	"sync"
	"time"
)

type Wait struct {
	waitMap sync.Map
	lock    sync.RWMutex
}

func NewWait() *Wait {
	return &Wait{}
}

func (w *Wait) Done(id uint32) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if v, ok := w.waitMap.Load(id); ok {
		info := v.(*Info)
		info.Ch <- struct{}{}
	}
}

// Wait 等待done调用，超时返回true
func (w *Wait) Wait(id, waitTime uint32, data interface{}) bool {
	ch := make(chan struct{})
	defer func() {
		w.lock.Lock()
		close(ch)
		w.waitMap.Delete(id)
		w.lock.Unlock()
	}()

	w.waitMap.Store(id, &Info{
		Ch:   ch,
		Data: data,
	})

	select {
	case <-time.After(time.Millisecond * time.Duration(waitTime)):
		return false
	case <-ch:
		return true
	}
}

// Data 获取缓存的数据
func (w *Wait) Data(id uint32) (interface{}, bool) {
	v, ok := w.waitMap.Load(id)
	if !ok {
		return nil, false
	}
	info := v.(*Info)
	return info.Data, true
}
