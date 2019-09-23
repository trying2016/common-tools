package utils

type QueueFun func(data interface{})

type Queue struct {
	queue      chan interface{}
	callBack   func(data interface{})
	exitSignal chan struct{}
}

func NewQueue(size int, fun func(data interface{})) *Queue {
	q := &Queue{}
	q.init(size, fun)
	return q
}

func (q *Queue) init(size int, fun func(data interface{})) {
	q.callBack = fun
	q.exitSignal = make(chan struct{})
	q.queue = make(chan interface{}, size)
	go q.run()
}

func (q *Queue) Push(value interface{}) {
	q.queue <- value
}

func (q *Queue) Destroy() {
	q.exitSignal <- struct{}{}
}

func (q *Queue) run() {
	for {
		select {
		case data := <-q.queue:
			q.callBack(data)
		case <-q.exitSignal:
			return
		}
	}
}
