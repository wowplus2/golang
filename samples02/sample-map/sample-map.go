package main

import "fmt"

var idMap map[int]string	// var 맵이름 map[Key타입]Value타입


func main() {
	fmt.Println("------------------------")

	idMap = make(map[int]string)	// map 초기화
	idMap[901] = "Apple"
	idMap[134] = "Grape"
	idMap[777] = "Tomato"

	// Key에 대한 값읽기
	str := idMap[134]
	fmt.Println("str :", str)

	noDate := idMap[999]	// 값이 없으면 nil or zero 리턴
	fmt.Println("noData :", noDate)

	// delete
	//delete(idMap, 777)

	for k,v := range(idMap) {
		fmt.Println(k, ":\t", v)
	}
	fmt.Println("------------------------")

	// 리터럴을 사용한 map 초기화
	tickers := map[string]string{
		"GOOG" :	"Google Inc",
		"MSFT" :	"Microsoft",
		"FB" :		"FaceBook",
	}

	for k,v := range(tickers) {
		fmt.Println(k, ":\t", v)
	}

	// Map의 Key 체크
	fmt.Println("I'll find 'TWET' in ticker map...")
	val, exists := tickers["TWET"]
	if !exists {
		fmt.Println("No TWET in ticker map!")
	} else {
		fmt.Println("'",val, "' in ticker Map")
	}
	fmt.Println("------------------------")
}
