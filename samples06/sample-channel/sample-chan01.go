package main

/*
[ Go 채널 ]
데이타를 주고 받는 통로라 볼 수 있는데, 채널은 make() 함수를 통해 미리 생성되어야 하며,
채널 연산자 <- 을 통해 데이타를 보내고 받는다.
채널은 흔히 goroutine들 사이 데이타를 주고 받는데 사용되는데, 상대편이 준비될 때까지 채널에서 대기함으로써
별도의 lock을 걸지 않고 데이타를 동기화하는데 사용된다.
*/

func main() {
	// 정수형 채널을 생성한다.
	ch := make(chan int)

	go func() {
		println("integer value send to channel...")
		ch <- 123	// 정수 123을 채널 ch에 전송한다.
	}()

	var i int
	i = <- ch	// 채널 ch로부터 정수 123을 수신한다.

	println(i)
}
