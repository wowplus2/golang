package main

import "fmt"

func main() {
	var nums []int = []int{5,4,3,2,1,0}
	fmt.Println("Unsorted: ", nums)

	bubbleSort(nums)
	fmt.Println("Sorted: ", nums)
}

func bubbleSort(nums []int) {
	var N int = len(nums)
	var i int

	for i = 0; i < N; i++ {
		if sweep(nums, i) {
			return
		}
	}
}

func sweep(nums []int, prevPasses int) bool {
	var N int = len(nums)
	var fIdx int = 0
	var sIdx int = 1
	var didSwap bool = false

	for sIdx < (N - prevPasses) {
		var fNum int = nums[fIdx]
		var sNum int = nums[sIdx]

		if fNum > sNum {
			nums[fIdx] = sNum
			nums[sIdx] = fNum
			didSwap = true
		}

		fIdx++
		sIdx++
	}

	return didSwap
}