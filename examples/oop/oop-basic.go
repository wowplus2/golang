package main

import (
	"math"
	"fmt"
)

type Vertx struct {
	X, Y float64
}

func (v *Vertx) Abs() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v *Vertx) Max() float64 {
	return math.Max(v.X, v.Y)
}


func main() {
	v := &Vertx{3, 4}
	fmt.Println("oop-sample01:")
	fmt.Println("\tmath.Abs() : ", v.Abs())
	fmt.Println("\tmath.Max() :", v.Max())
}
