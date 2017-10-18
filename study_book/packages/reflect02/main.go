// reflect - struct tag 출력
// 			: struct 정합성 규칙이나 JSON/XML 변환을 위한 정보를 표현
package main

import (
	"fmt"
	"reflect"
)

func main() {
	type User struct {
		Name string "check:len(3, 40)"
		Id   int    "check:range(1, 999999)"
	}

	u := User{"Myung", 1}
	uType := reflect.TypeOf(u)

	if fName, ok := uType.FieldByName("Name"); ok {
		fmt.Println(fName.Type, fName.Name, fName.Tag)
	}

	if fId, ok := uType.FieldByName("Id"); ok {
		fmt.Println(fId.Type, fId.Name, fId.Tag)
	}
}
