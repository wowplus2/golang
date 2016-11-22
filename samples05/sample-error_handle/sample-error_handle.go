package main

import (
	"os"
	"log"
)

func otherFunc() {

}

/*
func errorCheck() {
	_, err := otherFunc()

	switch err.(type) {
	default:	// no error
		println("ok")
	case MyError:
		log.Println("Log my error")
	case error:
		log.Fatal(err.Error())
	}
}
*/

func main() {
	//f, err := os.Open("C:\\Temp\\HncDownload\\1.txt")
	f, err := os.Open("C:\\Temp\\HncDownload\\Update.log")

	if err != nil {
		log.Fatal(err.Error())
	}

	println(f.Name())
}
