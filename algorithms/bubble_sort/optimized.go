package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var nums []int = rand.Perm(10)

	fmt.Println("Unsorted: ", nums)
	bubbledSort(nums)
	fmt.Println("Sorted: ", nums)
}

func bubbledUp(nums []int, n int) bool {
	var swapped bool = false
	for i := 0; i < n-1; i++ {
		var f int = nums[i]
		var s int = nums[i+1]
		if f > s {
			nums[i] = s
			nums[i+1] = f
			swapped = true
		}
	}

	return swapped
}

func bubbledSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		// Terminate early if bubbledUp doesn't swap anything, and reduce the max index we look at by len(nums) - 1
		if !bubbledUp(nums, len(nums) - 1) {
			return
		}
	}
}
