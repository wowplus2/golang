package main

func nval() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	next := nval()

	println(next())	// 1
	println(next())	// 2
	println(next())	// 3

	anothern := nval()
	println(anothern())	// 1 <- 다시시작
	println(anothern())	// 2
}
