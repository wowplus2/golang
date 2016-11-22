package main

import (
	"sync"
	"fmt"
)

func main() {
	// WaitGroup 생성 -> 2개의 goroutine을 기다림.
	var wait sync.WaitGroup
	wait.Add(2)

	// 익명함수를 사용한 goroutine
	go func() {
		defer wait.Done()	// 끝나면, .Done() 호출
		fmt.Println("Hello~")
	}()

	// 익명함수에 파라메터 전달
	go func(msg string) {
		defer wait.Done()	// 끝나면, .Done() 호출
		fmt.Println(msg)
	}("Hi!")

	wait.Wait()	// goroutine이 모두 끝날 때까지 대기...
}