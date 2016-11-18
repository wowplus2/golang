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

func sub_slice() {
	s := []int{0,1,2,3,4,5}
	s = s[2:5]

	fmt.Println("sub_slice: ")
	fmt.Println("\t", s)	// 2,3,4 출력
}

func append_slice() {
	s := []int{0, 1}
	// 하나의 요소 확장
	s = append(s, 2)	// 0, 1, 2
	// 여러개 요소 확장
	s = append(s, 3, 4, 5)	// 0, 1, 2, 3, 4, 5

	fmt.Println("append_slice: ")
	fmt.Println("\t", s)
}

func append_extend_slice1() {
	// len = 0, cap = 3 인 slice
	sliceA := make([]int, 0, 3)
	// 계속 한 요소씩 추가한다.
	fmt.Println("append_extend_slice1: ")
	for i := 1; i <= 15; i++ {
		sliceA = append(sliceA, i)
		// slice 길이와 용량 확인
		fmt.Println("\tlen =",len(sliceA),", cap =",cap(sliceA))
	}
	fmt.Println(sliceA)
}

func append_extend_slice2() {
	sliceA := []int{1,2,3}
	sliceB := []int{4,5,6}

	sliceA = append(sliceA, sliceB...)
	// sliceA = append(sliceA, 4,5,6)

	fmt.Println("append_extend_slice2: ")
	fmt.Println("\t",sliceA)
}

func copy_slice() {
	src := []int{0,1,2}
	trg := make([]int, len(src), cap(src)*2)

	copy(trg, src)

	fmt.Println("copy_slice: ")
	fmt.Println("\t",trg)	// [0 1 2]
	fmt.Println("\tlen =",len(trg),", cap =",cap(trg))
}


func main() {
	slice_way1()
	fmt.Println()
	slice_way2()
	fmt.Println()
	slice_way3()
	fmt.Println()
	sub_slice()
	fmt.Println()
	append_slice()
	fmt.Println()
	append_extend_slice1()
	fmt.Println()
	append_extend_slice2()
	fmt.Println()
	copy_slice()
}
