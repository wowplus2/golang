package main

import "fmt"

func main() {
	c := make(chan int, 2)

	c <- 1
	c <- 2
	go func() { c <- 3 }()
	go func() { c <- 4 }()
	go func() { c <- 5 }()

	fmt.Println("c <- 1 :", <-c)
	fmt.Println("c <- 2 :", <-c)
	fmt.Println("c <- 3 :", <-c)
	fmt.Println("c <- 4 :", <-c)
	fmt.Println("c <- 5 :", <-c)
}
