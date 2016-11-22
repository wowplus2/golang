package main

import "runtime"

/*
[ 다중 CPU 처리 ]
Go는 디폴트로 1개의 CPU를 사용한다.
즉, 여러 개의 Go 루틴을 만들더라도, 1개의 CPU에서 작업을 시분할하여 처리한다 (Concurrent 처리).
만약 머신이 복수개의 CPU를 가진 경우, Go 프로그램을 다중 CPU에서 병렬처리 (Parallel 처리)하게 할 수 있는데,
병렬처리를 위해서는 아래와 같이 runtime.GOMAXPROCS(CPU수) 함수를 호출하여야 한다.
(여기서 CPU 수는 Logical CPU 수를 가리킨다).
*/

func main() {
	// 4개의 CPU 사용...
	runtime.GOMAXPROCS(4)

	// ...
}
