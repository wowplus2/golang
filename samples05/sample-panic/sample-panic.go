package main

import "os"

// 현재 함수를 즉시 멈추고 현재 함수에 defer 함수들을 모두 실행한 후 즉시 리턴한다.
// 이러한 panic 모드 실행 방식은 다시 상위함수에도 똑같이 적용되고, 계속 콜스택을 타고 올라가며 적용된다.
// 그리고 마지막에는 프로그램이 에러를 내고 종료하게 된다.

func openFile(fn string) {
	f, err := os.Open(fn)

	if err != nil {
		panic(err)
	}
	// 파일 close 실행됨
	f.Close()
}

func main() {
	openFile("Invalid.txt")
	println("Done!")	// do not run this statement...
}