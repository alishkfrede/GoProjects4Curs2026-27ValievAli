package main

import "fmt"

type limiter interface {
	canAccept(n int) bool
	accept(n int) // сообщить, что приняли
}

// не больше N чисел
type countLimiter struct {
	max   int
	count int
}

func (l *countLimiter) canAccept(n int) bool { return l.count < l.max }
func (l *countLimiter) accept(n int)         { l.count++ }

// сумма не больше S
type sumLimiter struct {
	max int
	sum int
}

func (l *sumLimiter) canAccept(n int) bool { return l.sum+n <= l.max }
func (l *sumLimiter) accept(n int)         { l.sum += n }

func acceptAll(nums []int, l limiter) []int {
	res := make([]int, 0)
	for _, n := range nums {
		if l.canAccept(n) {
			l.accept(n)
			res = append(res, n)
		}
	}
	return res
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println("исходный:", nums)

	cl := &countLimiter{max: 3}
	fmt.Println("первые 3:", acceptAll(nums, cl))

	sl := &sumLimiter{max: 10}
	fmt.Println("сумма <=10:", acceptAll(nums, sl))
}