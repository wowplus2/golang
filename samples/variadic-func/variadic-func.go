package main


func say(msg ...string) {
	for _, s := range msg {
		println(s)
	}
}

func main() {
	say("This", "is", "a", "ebook")
	say("Hi!")
}
