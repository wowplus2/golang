package main

import "fmt"

func main() {
	// 빈 인터페이스는 어떠한 타입도 담을 수 있는 컨테이너라고 볼 수 있으며,
	// 여러 다른 언어에서 흔히 일컫는 Dynamic Type 이라고 볼 수 있다.
	var x interface{}
	x = 1
	x = "Tom"

	printIt(x)
}

func printIt(v interface{}) {
	fmt.Println(v)	// Tom
}
