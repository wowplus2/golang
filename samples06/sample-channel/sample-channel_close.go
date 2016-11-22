package main

func main() {
	ch := make(chan int, 2)

	// channel에 송신
	ch <- 100
	ch <- 200

	// channel을 닫는다.
	close(ch)

	// channel 수신
	println(<- ch)
	println(<- ch)

	if _, success := <- ch; !success {
		println("더 이상 데이터가 없습니다...")
	}
}
