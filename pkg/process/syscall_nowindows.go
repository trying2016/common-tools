//go:build !windows
// +build !windows

package process

import "syscall"

func NewSysProcAttr(hide bool) *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func SetConsoleMode() {
}

func GetLastError() error {
	return nil
}