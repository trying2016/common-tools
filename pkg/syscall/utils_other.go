//go:build !windows
// +build !windows

package util

import (
	"fmt"
	"syscall"
)

func SetConsoleMode() {

}

func SetMaxStdio() {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println(rLimit)

	rLimit.Max = 999999
	rLimit.Cur = 999999

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	fmt.Println("Rlimit Final", rLimit)
}

func LoadShareLib(name string) {

}

func EnableLinkedConnections() {

}
