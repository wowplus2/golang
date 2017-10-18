// reflect - 타입/값 정보 확인
package main

import (
	"fmt"
	ref "reflect"
)

func main() {
	x := 1
	y := 1.1
	z := "one"

	fmt.Printf("x: %v\t(%v)\n", ref.ValueOf(x).Int(), ref.TypeOf(x))
	fmt.Printf("x: %v\t(%v)\n", ref.ValueOf(y).Float(), ref.TypeOf(y))
	fmt.Printf("x: %v\t(%v)\n", ref.ValueOf(z).String(), ref.TypeOf(z))
}
