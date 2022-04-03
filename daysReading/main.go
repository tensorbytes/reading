package main

import (
	"fmt"
)

func main() {
	a := []int{1, 1, 5, 5, 8, 8, 4, 2, 3, 2, 3}
	fmt.Println(FindOddTime(a))
	b := 4
	c := 5
	Swap(b, c)
}

func Mean(max, min int32) int32 {
	return (max-min)>>2 + min
}

func FindOddTime(all []int) int {
	var sumxor int
	for _, n := range all {
		sumxor = sumxor ^ n
	}
	return sumxor
}

func Swap(a, b int) {
	a = a ^ b
	b = a ^ b
	a = a ^ b
}

// func FindEvenTime(all []int) int {
// 	return
// }
