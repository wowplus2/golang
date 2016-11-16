package main

import "fmt"

func main() {
	var a []int	// 슬라이스 변수 선언
	a = []int{1,2,3}
	a[1] = 10

	fmt.Println(a)	// [1,10,3] 출력
}
