package bunchRun

import (
	"reflect"
	"sync"
	"sync/atomic"
)

func Run(list interface{}, thread int, run func(v interface{}) bool) {
	listValue := reflect.ValueOf(list)
	var index = int64(0)
	count := int64(listValue.Len())
	var job sync.WaitGroup
	job.Add(thread)

	for i := 0; i < thread; i++ {
		go func() {
			defer job.Done()
			for {
				curIndex := atomic.AddInt64(&index, 1) - 1
				if curIndex >= count {
					return
				}
				elem := listValue.Index(int(curIndex))
				if ok := run(elem.Interface()); !ok {
					return
				}
			}
		}()
	}
	job.Wait()
}
