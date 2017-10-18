// reflect - 함수/메서드 동적 호출2
package main

import (
	"container/list"
	"fmt"
	ref "reflect"
)

func Len(x interface{}) int {
	value := ref.ValueOf(x)

	switch ref.TypeOf(x).Kind() {
	case ref.Array, ref.Chan, ref.Map, ref.Slice, ref.String:
		return value.Len()
	default:
		if method := value.MethodByName("Len"); method.IsValid() {
			values := method.Call(nil)
			return int(values[0].Int())
		}
	}

	panic(fmt.Sprintf("'%v' does not have a length", x))
}

func main() {
	a := list.New()
	b := list.New()
	b.PushFront(0.5)

	c := map[string]int{"A": 1, "B": 2}
	d := "one"
	e := []int{5, 0, 4, 1}

	fmt.Println(Len(a), Len(b), Len(c), Len(d), Len(e))
}
