package main

import (
	"fmt"
	"time"
)

func long() {
	fmt.Println("long 함수 시작.\t", time.Now())
	time.Sleep(3 * time.Second)
	fmt.Println("long 함수 종료.\t", time.Now())
}

func short() {
	fmt.Println("short 함수 시작.", time.Now())
	time.Sleep(1 * time.Second)
	fmt.Println("short 함수 종료.", time.Now())
}

func main() {
	fmt.Println("main 함수 시작.\t", time.Now())

	go long()
	go short()

	time.Sleep(5 * time.Second)
	fmt.Println("main 함수 종료.\t", time.Now())
}
