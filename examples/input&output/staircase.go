package main

import (
	"bufio"
	"os"
	"strings"
	"fmt"
	"strconv"
)

func main() {
	ir := bufio.NewReader(os.Stdin)
	N, _ := ir.ReadString('\n')
	num, _ := strconv.Atoi(strings.TrimSpace(N))
	tmp := num

	for i := tmp; i > 0; i-- {
		res := fmt.Sprintf("%*s", tmp, "#")
		for j := tmp; j < num; j++ {
			res += "#"
		}

		fmt.Println(res)
		tmp--
	}
}
