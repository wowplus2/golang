/*
	cf. try-catch
*/
package main

import (
	"fmt"
	"log"
)

func protect(g func()) {
	defer func() {
		log.Println("done")

		if err := recover(); err != nil {
			log.Printf("run time panic: %v", err)
		}
	}()

	log.Println("start")
	g()
}

func devide(a, b int) int {
	return a / b
}

func main() {
	protect(func() {
		fmt.Println(devide(4, 0))
	})
}
