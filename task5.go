package main

import "fmt"

type aggregator interface {
	add(n int)
	finish()
	result() string
}

// сумма и количество
type sumCountAgg struct {
	sum   int
	count int
}

func (a *sumCountAgg) add(n int) { a.sum += n; a.count++ }
func (a *sumCountAgg) finish()   {}
func (a *sumCountAgg) result() string {
	return fmt.Sprintf("sum=%d, count=%d", a.sum, a.count)
}

// минимум и максимум
type minMaxAgg struct {
	min, max int
	set      bool
}

func (a *minMaxAgg) add(n int) {
	if !a.set {
		a.min, a.max, a.set = n, n, true
		return
	}
	if n < a.min {
		a.min = n
	}
	if n > a.max {
		a.max = n
	}
}
func (a *minMaxAgg) finish() {}
func (a *minMaxAgg) result() string {
	if !a.set {
		return "empty"
	}
	return fmt.Sprintf("min=%d, max=%d", a.min, a.max)
}

func aggregate(nums []int, a aggregator) string {
	for _, n := range nums {
		a.add(n)
	}
	a.finish()
	return a.result()
}

func main() {
	nums := []int{-5, 3, -1, 10, 0}
	fmt.Println("исходный:", aggregate(nums, &sumCountAgg{}))
	fmt.Println("min/max:", aggregate(nums, &minMaxAgg{}))

	empty := []int{}
	fmt.Println("пустой sum/count:", aggregate(empty, &sumCountAgg{}))
	fmt.Println("пустой min/max:", aggregate(empty, &minMaxAgg{}))
}