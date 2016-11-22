package main

import "fmt"

// Rect - struct 정의
type Rect struct {
	width, height int
}

// Rect의 area() 메소드
//     Value Receiver( area메소드 내에서 struct 값이 변경되더라도 호출자의 데이터는 변경되지 않는다. )
func (r Rect) area() int {
	return r.width * r.height
}

// Rect의 area2() 메소드
//     Pointer Receiver( area2메소드 내의 값이 변경되면 변경된 값이 그대로 호출자에게 반영된다. )
func (r *Rect) area2() int {
	r.width++
	return r.width * r.height
}


func main() {
	rect := Rect{10, 20}
	fmt.Println("rect :", rect)

	area := rect.area()
	area2 := rect.area2()
	fmt.Println("[Value Receiver] 사각형의 면적(w:",rect.width,",h:",rect.height,") :", area)
	fmt.Println("[Pointer Receiver] 사각형의 면적(w:",rect.width,",h:",rect.height,") :", area2)
}
