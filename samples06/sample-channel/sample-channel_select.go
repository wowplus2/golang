package main

import "time"

// [ Channel Select ]
/*
select문은 복수 채널들을 기다리면서 준비된 (데이타를 보내온) 채널을 실행하는 기능을 제공한다.
즉, select문은 여러 개의 case문에서 각각 다른 채널을 기다리다가 준비가 된 채널 case를 실행하는 것이다.
select문은 case 채널들이 준비되지 않으면 계속 대기하게 되고, 가장 먼저 도착한 채널의 case를 실행한다.
만약 복수 채널에 신호가 오면, Go 런타임이 랜덤하게 그 중 한 개를 선택한다.
하지만, select문에 default 문이 있으면, case문 채널이 준비되지 않더라도 계속 대기하지 않고 바로 default문을 실행한다.

아래 예제는 for 루프 안에 select 문을 쓰면서 두개의 goroutine이 모두 실행되기를 기다리고 있다.
첫번째 run1()이 1초간 실행되고 done1 채널로부터 수신하여 해당 case를 실행하고, 다시 for 루프를 돈다.
for루프를 다시 돌면서 다시 select문이 실행되는데, 다음 run2()가 2초후에 실행되고 done2 채널로부터 수신하여 해당 case를 실행하게 된다.
done2 채널 case문에 break EXIT 이 있는데, 이 문장으로 인해 for 루프를 빠져나와 EXIT 레이블로 이동하게 된다.

Go의 "break 레이블" 문은 C/C# 등의 언어에서의 goto 문과 다른데,
Go에서는 해당 레이블로 이동한 후 자신이 빠져나온 루프 다음 문장을 실행하게 된다.
따라서, 여기서는 for 루프 다음 즉 main() 함수의 끝에 다다르게 된다.
*/
func main() {
	done1 := make(chan bool)
	done2 := make(chan bool)

	go run1(done1)
	go run2(done2)

	EXIT:
	for {
		select {
		case <- done1:
			println("run1 완료...")
		case <- done2:
			println("run2 완료...")
			break EXIT
		}
	}
}

func run1(done chan bool) {
	time.Sleep(1 * time.Second)
	done <- true
}

func run2(done chan bool) {
	time.Sleep(2 * time.Second)
	done <- true
}