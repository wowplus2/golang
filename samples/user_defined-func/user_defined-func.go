package main

// 구조체(struct)나 interface처럼 Custom Type(User Defined Type) 의 원형정의
type calculator func(int, int) int

// calculator 원형사용
func calc(f calculator, a int, b int) int {
	res := f(a, b)
	return res
}
