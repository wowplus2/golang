package main

import "fmt"

/*
[ 채널 파라메터 ]
채널을 함수의 파라미터도 전달할 때,
일반적으로 송수신을 모두 하는 채널을 전달하지만, 특별히 해당 채널로 송신만 할 것인지 혹은 수신만할 것인지를 지정할 수도 있다.
송신 파라미터는 (p chan<- int)와 같이 chan<- 을 사용하고,
수신 파라미터는 (p <-chan int)와 같이 <-chan 을 사용한다.
만약 송신 채널 파라미터에서 수신을 한다거나, 수신 채널에 송신을 하게되면, 에러가 발생한다.
*/

func main() {
	ch := make(chan string, 1)

	sendChannel(ch)
	recvChannel(ch)
}

func sendChannel(ch chan <- string) {
	ch <- "Data"
	// x := <- ch	// <==== error occurred!!
}

func recvChannel(ch <- chan string) {
	data := <- ch
	fmt.Println(data)
}
