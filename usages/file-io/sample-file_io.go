package main

import (
	"os"
	"io"
)

func main() {
	// 입력파일 열기
	fi, err := os.Open("C:\\Temp\\HncDownload\\Update.log")
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// 출력파일 생성
	fo, err := os.Create("C:\\Temp\\HncDownload\\Create.log")
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	buff := make([]byte, 1024)

	// loop
	for {
		// read
		cnt, err := fi.Read(buff)
		if err != nil && err != io.EOF {
			panic(err)
		}
		// 끝이면 loop 종료
		if cnt == 0 {
			break
		}

		// write
		_, err = fo.Write(buff[:cnt])
		if err != nil {
			panic(err)
		}
	}
}
