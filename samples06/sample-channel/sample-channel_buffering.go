package main

import "fmt"

// Dead Lock 발생 예시 : Unbuffered Channel
/*
func main() {
	c := make(chan int)

	c <- 1	// 수신루틴이 없으므로 dead lock 발생
	fmt.Println(<- c)	// 별도의 goroutine이 없으므로 comment해도 dead lock 발생
}
*/

// Buffered Channel
func main() {
	ch := make(chan int, 1)

	// 수신자가 없더라도 보낼 수 있다.
	ch <- 101

	fmt.Println(<- ch)
}