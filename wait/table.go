package wait

type Info struct {
	Ch   chan struct{}
	Data interface{}
}
