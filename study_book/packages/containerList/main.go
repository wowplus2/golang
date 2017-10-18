package main

import (
	"container/list"
	"fmt"
	"strings"
)

func main() {
	items := list.New()

	for _, x := range strings.Split("ABCDEFGH", "") {
		items.PushBack(x)
	}

	e := items.PushFront(0)
	items.InsertAfter(1, e)

	for el := items.Front(); el != nil; el = el.Next() {
		fmt.Printf("%v ", el.Value)
	}
}
