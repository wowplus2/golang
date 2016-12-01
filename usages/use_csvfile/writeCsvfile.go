package main

import (
	"os"
	"log"
	"encoding/csv"
	"bufio"
)

func main() {
	// file 생성
	file, err := os.Create("C:\\Go\\assets\\go_output.csv")
	if err != nil {
		log.Println(err)
	}

	// csv writer 생성
	wr := csv.NewWriter(bufio.NewWriter(file))

	// csv 내용 쓰기
	wr.Write([]string{"A", "0.25"})
	wr.Write([]string{"B", "3.141592"})
	wr.Flush()
}
