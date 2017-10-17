package main

import (
	"fmt"
	"log"
)

type fType func(int, int) int

func errorHandler(fn fType) fType {
	return func(a int, b int) int {
		defer func() {
			if err, ok := recover().(error); ok {
				log.Printf("run time panic -> %v", err)
			}
		}()

		return fn(a, b)
	}
}

func devide(a int, b int) int {
	return a / b
}

func main() {
	fmt.Println(errorHandler(devide)(4, 2))
	fmt.Println(errorHandler(devide)(3, 0))
}
