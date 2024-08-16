// Package thread
/**
 * @Author: trying
 * @Description:
 * @File:  thread.go
 * @Version: 1.0.0
 * @Date: 2024/8/16 22:32
 */

// thread.go
package thread

/*
#include "thread.h"
#include <stdio.h>
void threadFunc(void *arg);
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// CallbackWrapper 定义一个结构体来保存闭包和状态
type CallbackWrapper struct {
	callback func()
}

//export threadFunc
func threadFunc(arg unsafe.Pointer) {
	v := *(*uintptr)(arg)
	fn := *(*func())(unsafe.Pointer(v))
	fn()
}

// SetAffinity 设置线程亲和性
func SetAffinity(cpu int) int {
	return int(C.set_thread_affinity(C.int(cpu)))
}

// CreateThread 创建线程
func CreateThread(threadId int, fn func()) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	// Allocate memory in C for the function pointer
	cFn := C.thread_func_t(C.threadFunc)
	arg := C.malloc(C.size_t(unsafe.Sizeof(unsafe.Pointer(&fn))))
	*(*uintptr)(arg) = uintptr(unsafe.Pointer(&fn))
	// Create the thread
	C.create_thread(C.uint(threadId), cFn, unsafe.Pointer(arg))
}

// HardwareConcurrency 获取硬件线程数
func HardwareConcurrency() int {
	return int(C.hardware_concurrency())
}
