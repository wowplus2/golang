package main

import (
	"log"
	"os"
)

var myLogger *log.Logger

func main() {
	myLogger = log.New(os.Stdout, "[INFO] ", log.LstdFlags)

	// ... something to do ...
	run()
	myLogger.Println("End of Program...")
}

func run() {
	myLogger.Print("Test...")
}