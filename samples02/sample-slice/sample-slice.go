package main

import "fmt"


func slice_way1() {
	var a []int	// 슬라이스 변수 선언
	a = []int{1,2,3}
	a[1] = 10

	fmt.Println("slice_way1: ")
	fmt.Println("\t",a)	// [1,10,3] 출력
}

func slice_way2() {
	s := make([]int, 5, 10)

	fmt.Println("slice_way2: ")
	fmt.Println("\tlen =", len(s), ", cap =", cap(s))	// len 5, cap 10
}

func slice_way3() {
	var s[]int

	fmt.Println("slice_way3: ")
	if s == nil {
		fmt.Println("\tNil Slice")
	}
	fmt.Println("\tlen =", len(s), ", cap =", cap(s))	// len 0, cap 0
}

func sub_slice1() {
	s := []int{0,1,2,3,4,5}
	s = s[2:5]

	fmt.Println("sub_slice1: ")
	fmt.Println("\t", s)	// 2,3,4 출력
}


func main() {
	slice_way1()
	fmt.Println()
	slice_way2()
	fmt.Println()
	slice_way3()
	fmt.Println()
	sub_slice1()
}
