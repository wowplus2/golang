// reflect - 값 변경
package main

import (
	"fmt"
	ref "reflect"
)

func main() {
	lang := []string{"golang", "java", "c#"}
	fmt.Println(lang)

	sval := ref.ValueOf(lang)
	val := sval.Index(1)
	val.SetString("python")
	fmt.Println(lang)

	x := 1
	if v := ref.ValueOf(x); v.CanSet() {
		v.SetInt(2)
	}

	fmt.Println(x)

	v := ref.ValueOf(&x)
	p := v.Elem()
	p.SetInt(3)

	fmt.Println(x)
}
