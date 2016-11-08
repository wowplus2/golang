package seq

func Fib(n int) int {
	p, q := 0, 1
	for i := 0; i < n; i++ {
		p, q = q, p + q
	}
	return p
}
