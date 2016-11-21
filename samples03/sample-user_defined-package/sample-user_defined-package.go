package sample_user_defined_package

import "fmt"


var pop map[string]string


func init() {
	pop = make(map[string]string)

	pop["Adele"]		= "Hello"
	pop["Alicia Keys"]	= "Fallin"
	pop["John Legend"]	= "All of Me"
}

// GetMusic : Popular music by singer (외부에서 호출 가능)
func GetMusic(singer string) string {
	return pop[singer]
}

// 내부에서만 호출 가능
func getKeys() {
	for _,vp := range pop {
		fmt.Println(vp)
	}
}
