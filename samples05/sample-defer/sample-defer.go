package main

import "os"

// defer 키워드는 특정 문장 혹은 함수를 나중에 (defer를 호출하는 함수가 리턴하기 직전에) 실행하게 한다.
// 일반적으로 defer는 C#, Java 같은 언어에서의 finally 블럭처럼 마지막에 Clean-up 작업을 위해 사용된다.

func main() {
	f, err := os.Open("C:\\Temp\\HncDownload\\Update.log")

	if err != nil {
		panic(err)
	}

	// main 마지막에 file close 실행
	defer f.Close()

	// file read
	bytes := make([]byte, 1024)
	f.Read(bytes)

	println(len(bytes))
}
