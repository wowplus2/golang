package main

import "fmt"

type entity float32

func (e *entity) inc() {
	*e++
}

func (e *entity) echo() {
	fmt.Println(*e)
}

func main() {
	var e entity = 3
	e.echo()
	e.inc()
	e.echo()
	e.inc()
	e.echo()
}
