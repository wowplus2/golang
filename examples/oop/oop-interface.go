package main

import (
	"math"
	"fmt"
)


// interface를 만든다.
// 이 interface는 Area 함수를 가지고 있다.
type Shaper interface {
	Area() int
}

type Rectangle struct {
	width, height int
}

type Triangle struct {
	width, height int
}

type Circle struct {
	radius float64
}

// Area 메써드 구현
func (r Rectangle) Area() int {
	return r.width * r.height
}

func (r Triangle) Area() int {
	return (r.width * r.height) / 2
}

func (r Circle) Area() float64 {
	return (r.radius * r.radius) * math.Pi
}


func main() {
	r := Rectangle{3, 5}
	t := Triangle{3, 6}
	c := Circle{10.0}

	fmt.Println("oop-interface:")
	fmt.Println("\tArea of the Rectangle:\t", r.Area())
	fmt.Println("\tArea of the Triangle:\t", t.Area())
	fmt.Println("\tArea of the Circle:\t", c.Area())
	fmt.Println()
	s := Shaper(r)
	fmt.Println("\tArea of the Shaper r is ", s.Area())
}