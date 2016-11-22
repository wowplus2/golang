package main

import (
	"math"
	"fmt"
)

// 인터페이스 정의 : method들의 집합체
// 	타입(type)이 구현해야하는 메서드 원형(prototype)들을 정의한다.
type Shape interface {
	area() float64
	perimeter() float64
}

// Rect sctuct 정의
type Rect struct {
	width, height float64
}

// Circle struct 정의
type Circle struct {
	radius float64
}

// Rect struct 타입에 대한 Shape 인터페이스 구현
func (r Rect) area() float64 { return r.width * r.height }
func (r Rect) perimeter() float64 {
	return 2 * (r.width + r.height)
}

// Circle struct 타입에 대한 Shape interface 구현
func (c Circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c Circle) perimeter() float64 {
	return 2 * math.Pi * c.radius
}


// 인터페이스 사용
func main() {
	r := Rect{10., 20.}
	c := Circle{10}

	printArea(r, c)
}

func printArea(shapes ...Shape) {
	for _, s := range shapes {
		a := s.area()	// interface method 호출
		fmt.Println(a)
	}
}