package main

import (
	"fmt"
	"os"
)

// panic 함수에 의한 패닉상태를 다시 정상상태로 되돌리는 함수이다.

func openFile(fn string) {
	// defer 함수. panic 호출 시 실행됨.
	defer func() {
		if r := recover(); r != nil {
			fmt.Print("OPEN ERROR : ", r)
		}
	}()

	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	// file close 실행됨
	defer f.Close()
}

func main() {
	openFile("1.txt")
	println("Done by recover method...!!")
}