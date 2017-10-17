package main

import (
	"fmt"
	"runtime"
)

func showOS() {
	fmt.Println("현재 파일: util_windows.go")
	fmt.Println(runtime.GOOS)
}
