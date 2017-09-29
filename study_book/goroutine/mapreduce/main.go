package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

const (
	BUF_SIZE = 10
)

var (
	path    = flag.String("path", "..", "find path")
	pattern = flag.String("pattern", ".go$", "grep pattern")
)

func find(dir, pattern string) <-chan string {
	out := make(chan string, 1000)

	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	go func() {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if regex.MatchString(path) {
				out <- path
			}

			return nil
		})
		close(out)
	}()

	return out
}

func parseArgs() (string, string) {
	flag.Parse()
	return *path, *pattern
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println(runConcurrentReduce(collect(runConcurrentMap(find(parseArgs())))))
}
