package main

import "fmt"

func main() {
	s := []string{"Daniel", "Myung"}
	// appending an element at the end of slice
	s = append(s, "Jina", "Jung")
	fmt.Println(s)
	// removing element at the end
	s = s[:len(s)-1]
	fmt.Println(s)
	idx := 1
	// removing element at the nth index
	s = append(s[:idx], s[idx+1:]...)
	fmt.Println(s)
}
