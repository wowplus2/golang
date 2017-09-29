package main

import "fmt"

type shaper interface {
	area() float64
}

type rect struct{ width, height float64 }

type circle struct{ radius float64 }

func (r rect) area() float64 {
	return r.width * r.height
}

func (r rect) show() {
	fmt.Printf("width: %f, height: %f\n", r.width, r.height)
}

func (c circle) show() {
	fmt.Printf("radius: %f", c.radius)
}

func describe(s shaper) {
	fmt.Println("area: ", s.area())
}

// 익명 인터페이스
func display(s interface{}) {
	//s.show()
	fmt.Println(s)
}

func main() {
	r := rect{3, 4}
	c := circle{2.5}

	describe(r)
	display(r)
	display(c)
	display(3.14)
	display("rect struct")
}
