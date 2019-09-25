package utils

import (
	"fmt"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	fmt.Printf("password:%s\n", generateHash("65f4b181", "653520"))
	if "sha1$5a7220be$1$9df5f7c2e133b96ca27d3b56181b2007a3662108" != generateHash("65f4b181", "653520") {
		t.Error("fail")
	}
}
