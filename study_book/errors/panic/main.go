package main

import "fmt"

func devide(a, b int) int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	return a / b
}

func main() {
	//fmt.Println("Starting the program")
	//panic("A severe error occured: stopping the program!")
	//fmt.Println("Ending the program")
	fmt.Println("result:", devide(1, 0))
}
