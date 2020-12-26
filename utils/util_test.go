package utils

import (
	"fmt"
	"testing"
)

func TestSubleString(t *testing.T) {
	str := "123abcd"
	fmt.Printf(SubleString(str, "3", "d"))
}

func TestSliceRemove(t *testing.T) {
	arr := [6]int{1, 2, 3, 4, 5, 6}
	SliceRemove(arr, 1)
}

func TestQueue(t *testing.T) {
	qu := NewQueue(100, func(data interface{}) {
		println(data)
	})
	qu.Push(1)
	qu.Push(1)
	count := qu.Count()
	println(count)
}

func TestToFloat(t *testing.T) {
	fValue := ToFloat("100.00")
	fmt.Printf("value %v", fValue)
}
