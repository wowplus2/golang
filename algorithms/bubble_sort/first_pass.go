package main

import (
	"math/rand"
	"fmt"
)

func main() {
	// First we are going to start by creating an array with the numbers 0 - 10 that are in a psuedo-random order.
	// These aren't truly random because I am not seeding the random package,
	// so they end up being the same every time we run our code, but they are not in sorted order.
	// This is actually really handy for cases like this since everyone gets to see the exact same execution.
	var nums []int = rand.Perm(10)

	// Next we are going to print out the numbers just to verify that they aren't sorted.
	// You should see the output: [9 4 2 6 8 0 3 1 7 5]
	fmt.Println("Unsorted: ", nums)

	// First we are going to focus on the most basic implementation possible.
	// That means we won't be doing any of the optimization we discussed in the article,
	// but instead will just be writing code that follows the basic algorithm described.
	//
	// To do this we are going to break it into two parts.
	// The first is the code that compares every pair in the array (slice in GO, but roughly the same) ans swaps
	// two values if the first is greater than the second.
	// This will be the function bubbleUp().
	bubbleSort(nums)
	fmt.Println("Sorted: ", nums)
}

// The purpose of this function is to compare every consecutive pair of numbers and swap them if the first is greater than
// the second. In essence it is bubbling up the largest number that isn't in its final position with each pass.
//
// NOTE: Slices (similar to arrays in other languages, but not exactly the same) are passed by reference in GO,
//       so any changes we make to this array will be propogated to the original slice.
// 	 You could just as easily return the array when this function finishes executing, but we won't be.
func bubbleUp(nums []int) {
	// len(nums) - 1 because we will look at the number at i and i+1, so we by setting our max as less than (not equal to)
	// len(nums) - 1 we ensure that we don't go out of bounds of our array.
	for i := 0; i < len(nums) - 1; i++ {
		// Get the pair of numbers that we want to look at.
		var first int = nums[i]
		var second int = nums[i + 1]
		// Check to see if the first number in the pair is greater than the second.
		if first > second {
			// If it is, we neet to swap the numbers.
			nums[i] = second
			nums[i + 1] = first
		}
	}
}

// The puepose of this function is sort the numbers slice by iterating N times (where N is the size of the slice)
// and calling the bubbleUp function on each iteration.
func bubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		bubbleUp(nums)
	}
}
