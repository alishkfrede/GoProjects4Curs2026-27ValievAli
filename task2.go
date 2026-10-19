package main

import "fmt"

// первая горутина: генерирует числа от 1 до N и пишет в канал
func generate(n int, out chan int) {
	for i := 1; i <= n; i++ {
		out <- i
	}
	close(out) // закрываем канал, когда числа закончились
}

// вторая горутина: читает числа, возводит в квадрат и пишет в другой канал
func square(in chan int, out chan int) {
	for num := range in { // range сам завершится, когда in закроется
		out <- num * num
	}
	close(out)
}

func main() {
	n := 10

	// два канала: для чисел и для квадратов
	numbers := make(chan int)
	squares := make(chan int)

	// запускаем горутины
	go generate(n, numbers)
	go square(numbers, squares)

	// main собирает квадраты в срез
	var result []int
	for sq := range squares {
		result = append(result, sq)
	}

	fmt.Println("Квадраты:", result)
}