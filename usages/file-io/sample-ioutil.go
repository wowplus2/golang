package main

import "io/ioutil"

func main() {
	// file read
	bytes, err := ioutil.ReadFile("C:\\Temp\\HncDownload\\Update.log")
	if err != nil {
		panic(err)
	}

	// file write
	err = ioutil.WriteFile("C:\\Temp\\HncDownload\\Create.ioutil.log", bytes, 0)
	if err != nil {
		panic(err)
	}
}
