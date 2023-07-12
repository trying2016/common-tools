//go:build windows
// +build windows

package util

/*
#include <windows.h>
#include <mmsystem.h>
#include <stdio.h>
void setConsoleMode(){
	HANDLE hInput = GetStdHandle(STD_INPUT_HANDLE);
	DWORD prev_mode;
	GetConsoleMode(hInput, &prev_mode);
	SetConsoleMode(hInput, ENABLE_EXTENDED_FLAGS | (prev_mode & ~ENABLE_QUICK_EDIT_MODE));
}

int SetMaxStdio(){
	int num = _getmaxstdio();
	_setmaxstdio(8192);
	return num;
}
*/
import "C"
import (
	reg "golang.org/x/sys/windows/registry"
	"runtime"
	"syscall"
)

const (
	FILE_CACHE_MAX_HARD_ENABLE = 0x1
)

func doGetProcAddress(lib uintptr, name string) uintptr {
	addr, _ := syscall.GetProcAddress(syscall.Handle(lib), name)
	return uintptr(addr)
}
func doLoadLibrary(name string) uintptr {
	lib, _ := syscall.LoadLibrary(name)
	return uintptr(lib)
}
func syscall3(trap, nargs, a1, a2, a3 uintptr) uintptr {
	ret, _, _ := syscall.Syscall(trap, nargs, a1, a2, a3)
	return ret
}

func SetSystemFileCacheSize(minimumFileCacheSize, maximumFileCacheSize uint64, flags int64) bool {
	libkernel32 := doLoadLibrary("kernel32.dll")
	setSystemFileCacheSize := doGetProcAddress(libkernel32, "SetSystemFileCacheSize")
	ret1 := syscall3(setSystemFileCacheSize, 3,
		uintptr(minimumFileCacheSize),
		uintptr(maximumFileCacheSize),
		uintptr(flags))
	return ret1 != 0
}

func SetConsoleMode() {
	C.setConsoleMode()
}

func SetMaxStdio() {
	C.SetMaxStdio()
}
func EnableLinkedConnections() {
	var KEY_ALL_ACCESS = 0x000F003F
	var KEY_WOW64_64KEY = 0x0100
	var access = KEY_ALL_ACCESS
	// 当前是32位的程序运行在x64下面
	//if win.IsWow64() {
	if runtime.GOARCH == "amd64" {
		access |= KEY_WOW64_64KEY
	}
	k, _, err := reg.CreateKey(reg.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System", uint32(access))
	if err != nil {
		return
	}
	defer k.Close()
	k.SetDWordValue("EnableLinkedConnections", uint32(1))
}
