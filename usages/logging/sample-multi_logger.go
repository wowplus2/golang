package main

import (
	"log"
	"os"
	"io"
)

func run() {
	log.Print("Test Multi logger...")
}

func main() {
	fpLog, err := os.OpenFile("multi_logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer fpLog.Close()

	// 파일과 cmd화면에 같이 출력하기 위해 MultiWriter 생성
	mWriter := io.MultiWriter(fpLog, os.Stdout)
	log.SetOutput(mWriter)

	run()

	log.Println("End of Milti new logger program...")
}
