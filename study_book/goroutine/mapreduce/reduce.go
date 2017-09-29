package main

import (
	"bytes"
	"fmt"
	"sync"
	"text/scanner"
)

type summary struct {
	// 키: token
	// 값: map[string]int
	// 			키: file path
	// 			값: token count
	m map[string]map[string]int

	// 공유 데이터 m을 보호하기 위한 뮤텍스
	mu sync.Mutex
}

func reducer(token string, positions []scanner.Position) map[string]int {
	result := make(map[string]int)
	for _, p := range positions {
		result[p.Filename] += 1
	}

	return result
}

func (s summary) String() string {
	var buff bytes.Buffer

	for token, val := range s.m {
		buff.WriteString(fmt.Sprintf("Token: %s\n", token))
		total := 0
		for path, cnt := range val {
			if path == "" {
				continue
			}
			total += cnt
			buff.WriteString(fmt.Sprintf("%8d %s", cnt, path))
			buff.WriteString("\n")
		}
		buff.WriteString(fmt.Sprintf("Total: %d\n\n", total))
	}

	return buff.String()
}

func runReduce(tokenPositions intermediate) summary {
	s := summary{m: make(map[string]map[string]int)}
	for token, positions := range tokenPositions {
		s.mu.Lock() //	m값을 변경하는 부분(임계 영역)을 뮤텍스로 잠금
		s.m[token] = reducer(token, positions)
		s.mu.Unlock() // m값 변경 완료 후 뮤텍스 잠금 해제
	}
	return s
}

func runConcurrentReduce(in intermediate) summary {
	s := summary{m: make(map[string]map[string]int)}
	var wg sync.WaitGroup

	for token, val := range in {
		wg.Add(1)
		go func(token string, positions []scanner.Position) {
			defer wg.Done()
			s.mu.Lock() // m 값을 변경하는 부분(임계 영역)을 뮤텍스로 잠금
			s.m[token] = reducer(token, positions)
			s.mu.Unlock() // m 값 변경 완료 후 뮤텍스 잠금 해제
		}(token, val)
	}

	wg.Wait()
	return s
}
