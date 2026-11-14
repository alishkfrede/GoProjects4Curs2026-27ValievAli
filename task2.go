package main

import "fmt"

type mapper interface {
	mapIt(n int) int
}

type doubleMapper struct{}

func (doubleMapper) mapIt(n int) int { return n * 2 }

type absMapper struct{}

func (absMapper) mapIt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func transform(nums []int, m mapper) []int {
	res := make([]int, len(nums))
	for i, n := range nums {
		res[i] = m.mapIt(n)
	}
	return res
}

func main() {
	nums := []int{-5, -1, 0, 3, -7}
	fmt.Println("исходный:", nums)
	fmt.Println("x2:", transform(nums, doubleMapper{}))
	fmt.Println("abs:", transform(nums, absMapper{}))
}