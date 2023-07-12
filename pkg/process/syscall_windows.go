//go:build windows
// +build windows

package process

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
	"golang.org/x/sys/windows"
	"syscall"
)

func NewSysProcAttr(hide bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{HideWindow: hide}
}

func SetConsoleMode() {
	C.setConsoleMode()
}

func GetLastError() error{
	return windows.GetLastError()
}