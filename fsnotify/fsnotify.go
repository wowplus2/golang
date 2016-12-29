package main

import (
	"github.com/howeyc/fsnotify"
	"log"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <- watcher.Event:
				log.Println("event :", ev)
			case err := <- watcher.Error:
				log.Println("error :", err)
			}
		}
	}()

	err = watcher.Watch("C:\\")
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit.
	<- done

	/* ... do stuff ... */
	watcher.Close()
}