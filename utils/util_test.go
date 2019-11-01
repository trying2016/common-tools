package utils

import (
	"fmt"
	"testing"
)

func TestSubleString(t *testing.T) {
	str:= "123abcd"
	fmt.Printf(SubleString(str, "3", "d"))
}