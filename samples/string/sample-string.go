package main

import "fmt"


func main() {
	// Raw String Literal. 복수라인.
	rLiteral := `아이랑\n
	아리랑\n
	아라리요~`

	// Interpreted String Literal
	iLiteral := "아리랑아리랑\n아라리요~"
	// 아래와 같이 +를 사용하여 두 라인에 걸쳐 사용할 수도 있다.
	// iLiteral := "아리랑아리랑\n" +
	//		"아라리요~"

	fmt.Println(rLiteral)
	fmt.Println()
	fmt.Println(iLiteral)
}
