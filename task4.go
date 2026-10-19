package main

import (
	"fmt"
	"sync"
)

// генератор 1: пишет свои числа в канал
func gen1(out chan int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

// генератор 2
func gen2(out chan int) {
	for i := 10; i <= 50; i += 10 {
		out <- i
	}
	close(out)
}

// генератор 3
func gen3(out chan int) {
	for i := 100; i <= 500; i += 100 {
		out <- i
	}
	close(out)
}

// fan-in: сливает несколько каналов в один
func merge(channels ...chan int) chan int {
	merged := make(chan int)

	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		// для каждого входного канала — своя горутина,
		// которая перекладывает числа в общий канал
		go func(c chan int) {
			defer wg.Done()
			for num := range c {
				merged <- num
			}
		}(ch)
	}

	// отдельная горутина: ждёт всех и закрывает merged
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	// три канала от трёх генераторов
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	// запускаем генераторы
	go gen1(ch1)
	go gen2(ch2)
	go gen3(ch3)

	// сливаем в один канал
	merged := merge(ch1, ch2, ch3)

	// собираем всё в срез
	var result []int
	for num := range merged {
		result = append(result, num)
	}

	fmt.Println("Все числа:", result)
}