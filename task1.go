package main

import "fmt"

type predicate interface {
	check(n int) bool
}

// чётное
type evenPred struct{}

func (evenPred) check(n int) bool { return n%2 == 0 }

// положительное
type positivePred struct{}

func (positivePred) check(n int) bool { return n > 0 }

func filter(nums []int, p predicate) []int {
	res := make([]int, 0)
	for _, n := range nums {
		if p.check(n) {
			res = append(res, n)
		}
	}
	return res
}

func main() {
	nums := []int{-3, -2, -1, 0, 1, 2, 3, 4, 5}
	fmt.Println("исходный:", nums)
	fmt.Println("чётные:", filter(nums, evenPred{}))
	fmt.Println("положительные:", filter(nums, positivePred{}))
}