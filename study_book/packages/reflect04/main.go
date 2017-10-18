// reflect - 함수/메서드 동적 호출
package main

import (
	"fmt"
	ref "reflect"
	"strings"
)

func TitleCase(s string) string {
	return strings.Title(s)
}

func main() {
	caption := "go is an open source programming language"
	// TitlaCase 바로 호출
	title := TitleCase(caption)
	fmt.Println(title)

	// TitleCase 동적 호출
	titleFuncVal := ref.ValueOf(TitleCase)
	vals := titleFuncVal.Call([]ref.Value{ref.ValueOf(caption)})
	title = vals[0].String()
	fmt.Println(title)
}
