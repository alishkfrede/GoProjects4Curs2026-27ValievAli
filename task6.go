package main

import "fmt"

type formatter interface {
	format(n int) string
}

type plainFmt struct{}

func (plainFmt) format(n int) string { return fmt.Sprintf("%d", n) }

type bracketFmt struct{}

func (bracketFmt) format(n int) string { return fmt.Sprintf("[%d]", n) }

type parityFmt struct{}

func (parityFmt) format(n int) string {
	tag := "нечёт"
	if n%2 == 0 {
		tag = "чёт"
	}
	return fmt.Sprintf("%d(%s)", n, tag)
}

func printAll(nums []int, f formatter) {
	for _, n := range nums {
		fmt.Print(f.format(n), " ")
	}
	fmt.Println()
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	printAll(nums, plainFmt{})
	printAll(nums, bracketFmt{})
	printAll(nums, parityFmt{})
}