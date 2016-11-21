package main

import "github.com/wowplus2/golang/samples03/sample-user_defined-package"

func main() {
	song := sample_user_defined_package.GetMusic("Alicia Keys")
	println(song)
}
