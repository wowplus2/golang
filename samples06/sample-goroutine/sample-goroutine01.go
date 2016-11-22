package main

import (
	"fmt"
	"time"
)

func say(s string) {
	for i := 0; i < 10; i++ {
		fmt.Println(s, "***", i)
	}
}

func main() {
	// 함수를 동기적으로 실행
	say("Sync Call")
	// 함수를 비동기적으로 실행
	go say("Async Call-01")
	go say("Async Call-02")
	go say("Async Call-03")

	// 3 seconds 대기
	time.Sleep(time.Second * 3)
}