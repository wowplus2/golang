package main

import (
	"log"
	"os"
)

func run() {
	log.Print("Test New logger...")
}

func main() {
	fpLog, err := os.OpenFile("logfile_new.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer fpLog.Close()

	// 표준로거를 파일로그로 변경
	log.SetOutput(fpLog)

	run()

	log.Println("End of New logger Program...")
}
