package main

func main() {
	msg := "Hello"
	say(&msg)
	println("main: " + msg)
}

func say(msg *string) {
	println("say: " + *msg)
	*msg = "Changed"
}
