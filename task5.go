package main

import (
	"context"
	"fmt"
	"time"
)

// "тяжёлая" обработка: умножаем на 10 и спим
func process(ctx context.Context, num int, results chan int) {
	select {
	case <-ctx.Done():
		// контекст отменён — выходим, ничего не делаем
		return
	case <-time.After(300 * time.Millisecond):
		// имитация тяжёлой работы
	}

	result := num * 10

	// если результат больше 50 — отменяем всю работу
	if result > 50 {
		fmt.Printf("Число %d дало результат %d > 50, отменяем!\n", num, result)
		// отмену делает main через cancel, а тут просто пишем результат
	}

	// пишем результат (но проверяем, не отменили ли нас)
	select {
	case results <- result:
	case <-ctx.Done():
		return
	}
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // на всякий случай

	results := make(chan int)

	// семафор на 2: буферизованный канал с 2 слотами
	sem := make(chan struct{}, 2)

	// запускаем по горутине на каждое число
	for _, num := range numbers {
		sem <- struct{}{} // занимаем слот (если 2 занято — ждём)

		go func(n int) {
			defer func() { <-sem }() // освобождаем слот
			process(ctx, n, results)
		}(num)
	}

	// собираем результаты в отдельной горутине
	var collected []int
	done := make(chan struct{})
	go func() {
		defer close(done)
		for res := range results {
			collected = append(collected, res)
			if res > 50 {
				cancel() // отменяем всё!
				return   // выходим из чтения
			}
		}
	}()

	<-done
	fmt.Println("Собрали:", collected)
}