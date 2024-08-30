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
	for i := 0; i < 16; i++ {
		CreateThread(i, func() {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Second)
				t.Log("CreateThread", i)
			}
		})
	}
	time.Sleep(time.Second * 10)
}

func TestHardwareConcurrency(t *testing.T) {
	t.Log("HardwareConcurrency", HardwareConcurrency(),
		"CPU(s)", runtime.NumCPU())
}
