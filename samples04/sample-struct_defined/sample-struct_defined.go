package main

import "fmt"

// struct 정의
type person struct {
	name	string
	age	int
}


func main() {
	// person object 생성
	p := person{}
	// other ways...
	var p1 person
	p1 = person{"Bob", 20}

	p2 := person{ name: "Daniel", age: 41 }

	// usage way of 'new'
	p3 := new(person) // 객체의 포인터를 리턴한다.
	p3.name = "Lee"	// p3가 포인터라도 . 을 사용한다.

	// field값 설정
	p.name	= "Myung"
	p.age	= 41

	fmt.Println("p: ", p)
	fmt.Println("p1: ", p1)
	fmt.Println("p2: ", p2)
	fmt.Println("p3: ", *p3)	// used pointer
}
