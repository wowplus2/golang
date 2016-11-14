package main

func main() {
	// 변수 add 에 익명함수 할당
	add := func(i int, j int) int {
		return i + j
	}

	// add 함수 전달
	r1 := calc(add, 10, 20)
	println(r1)

	// 직접 1번째 파라메터에 익명함수를 정의함
	r2 := calc(func(x int, y int) int { return x - y }, 10, 20)
	println(r2)
}

func calc(f func(int, int) int, a int, b int) int {
	res := f(a, b)
	return res
}
