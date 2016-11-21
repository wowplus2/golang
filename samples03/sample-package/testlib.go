package testlib

// 패키지안에 init()한수만 호출
//import _ "other/xlinb"

// 패키지 alias
/*
import (
	mongo "other/mongo/db"
	mysql "other/mysql/db"
)

func main() {
	mondb := mongo.Get()
	mydb := mysql.Get()
	// ...
}*/

var pop map[string]string

// package 로드 시 map 초기화
func init() {
	pop = make(map[string]string)
}
