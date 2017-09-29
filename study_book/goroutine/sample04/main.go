package main

import (
	"fmt"
	"time"
)

/*
func fibonacci(c, quit chan int) {
	x, y := 0, 1

	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func actFibonacci() {
	c := make(chan int)
	quit := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()

	fibonacci(c, quit)
}
*/

func actSelectDefaultCase() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)

	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!!")
			return
		default:
			fmt.Println("     .")
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func main() {
	//actFibonacci()
	//fmt.Println()
	actSelectDefaultCase()
}
