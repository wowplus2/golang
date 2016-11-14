package main

func main() {
	sum := func(n ...int) int {	// 익명함수 정의
		s := 0
		for _, i := range n {
			s += i
		}
		return s
	}

	res := sum(1,2,3,4,5,6,7,8,9,10) // 익명함수 호출
	println(res)
}
