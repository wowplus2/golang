package main


func say(msg ...string) {
	for _, s := range msg {
		println(s)
	}
}

// Named Return Parameters - 1
func sum(nums ...int) int {
	s := 0
	for _, n := range nums {
		s += n
	}
	return s
}

// Named Return Parameters - 2
/*func sum2(nums ...int) (int, int) {
	s := 0
	count := 0
	for _, n := range nums {
		s += n
		count++
	}
	return count, s
}*/
func sum2(nums ...int) (count int, total int) {
	for _, n := range nums {
		total += n
	}
	count = len(nums)
	return
}

func main() {
	total := sum(1,7,3,5,9)
	count, t := sum2(1,7,3,5,9)

	say("This", "is", "a", "ebook")
	say("Hi!")
	println(total)
	println(count, t)
}
