package main

import "fmt"

func main() {
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(i)
		}
		done <- true
	}()

	// 위의 goroutine이 끝날때까지 대기...
	<- done
}
