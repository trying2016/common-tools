package bunchRun

import "testing"

func TestRun(t *testing.T) {
	var list []int
	for i := 0; i < 10; i++ {
		list = append(list, i*2)
	}
	Run(list, 10, func(i interface{}) {
		t.Log(i)
	})
}
