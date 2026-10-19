package main

import (
	"fmt"
	"sync"
)

// воркер: читает числа из канала jobs, умножает на 2 и пишет в results
func worker(jobs chan int, results chan int, wg *sync.WaitGroup) {
	defer wg.Done() // сообщаем, что воркер закончил
	for num := range jobs {
		results <- num * 2
	}
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	jobs := make(chan int)
	results := make(chan int)

	// запускаем 3 воркера
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// отдельная горутина кидает числа в jobs и закрывает канал
	go func() {
		for _, num := range numbers {
			jobs <- num
		}
		close(jobs)
	}()

	// отдельная горутина ждёт всех воркеров и закрывает results
	go func() {
		wg.Wait()
		close(results)
	}()

	// main собирает результаты
	var output []int
	for res := range results {
		output = append(output, res)
	}

	fmt.Println("Результаты:", output)
}