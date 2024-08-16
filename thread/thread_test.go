// Package thread
/**
 * @Author: trying
 * @Description:
 * @File:  thread_test.go
 * @Version: 1.0.0
 * @Date: 2024/8/16 22:34
 */

package thread

import (
	"runtime"
	"testing"
	"time"
)

func TestSetAffinity(t *testing.T) {
	t.Log(SetAffinity(1))
}

func TestCreateThread(t *testing.T) {
	CreateThread(1, func() {
		t.Log("Hello, thread!")
	})
	time.Sleep(time.Second)
}

func TestHardwareConcurrency(t *testing.T) {
	t.Log("HardwareConcurrency", HardwareConcurrency(),
		"CPU(s)", runtime.NumCPU())
}
