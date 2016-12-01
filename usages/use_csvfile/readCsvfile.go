package main

import (
	"os"
	"encoding/csv"
	"bufio"
	"fmt"
)

func main() {
	// 파일 open
	file, _ := os.Open("C:\\Go\\assets\\g4_member.csv")

	// csv reader 생성
	rdr := csv.NewReader(bufio.NewReader(file))

	// csv 내용 모두읽기
	rows, _ := rdr.ReadAll()

	// 행, 열 읽기
	for i, row := range rows {
		for j := range row {
			fmt.Printf("%s\t", rows[i][j])
		}
		fmt.Println()
	}
}
