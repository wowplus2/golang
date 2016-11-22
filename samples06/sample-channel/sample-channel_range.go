package main


// 채널에서 송신자가 송신을 한 후, 채널을 닫을 수 있다. 그리고 수신자는 임의의 갯수의 데이타를 채널이 닫힐 때까지 계속 수신할 수 있다.
// 채널 range문은 range 키워드 다음의 채널로부터 계속 수신하다가 채널이 닫힌 것을 감지하면 for 루프를 종료한다.
func main() {
	ch := make(chan int, 2)

	// channel에 송신
	ch <- 100
	ch <- 200

	// channel을 닫는다.
	close(ch)

	// 방법.1
	//   channel이 close될 때까지 계속 수신
/*
	for {
		if i, success := <- ch; success {
			println(i)
		} else {
			break
		}
	}
*/
	// 방법.2
	//  방법.1 과 동일한 channel range 문
	for i := range ch {
		println(i)
	}
}
