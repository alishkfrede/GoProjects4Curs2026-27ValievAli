package main

import "fmt"

// функция считает сумму чисел в переданном ей куске среза
// и отправляет результат в канал
func sumPart(numbers []int, result chan int) {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	result <- sum
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// делим срез пополам
	middle := len(numbers) / 2
	firstHalf := numbers[:middle]
	secondHalf := numbers[middle:]

	// создаём канал, куда горутины положат результаты
	result := make(chan int)

	// запускаем две горутины
	go sumPart(firstHalf, result)
	go sumPart(secondHalf, result)

	// забираем два результата из канала и складываем
	sum1 := <-result
	sum2 := <-result
	total := sum1 + sum2

	fmt.Println("Сумма первой половины:", sum1)
	fmt.Println("Сумма второй половины:", sum2)
	fmt.Println("Общая сумма:", total)
}