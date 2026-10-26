package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// ---------- утилита ----------

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ==========================================================
// ЗАДАЧА 1. Сумма среза двумя горутинами
// ==========================================================

func sumPart(numbers []int, result chan int) {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	result <- sum
}

type SumRequest struct {
	Numbers []int `json:"numbers"`
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req SumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	middle := len(req.Numbers) / 2
	result := make(chan int, 2)

	go sumPart(req.Numbers[:middle], result)
	go sumPart(req.Numbers[middle:], result)

	sum1 := <-result
	sum2 := <-result

	writeJSON(w, map[string]int{
		"sum1":  sum1,
		"sum2":  sum2,
		"total": sum1 + sum2,
	})
}

// ==========================================================
// ЗАДАЧА 2. Pipeline: генератор → квадратор → main
// ==========================================================

func generate(n int, out chan int) {
	for i := 1; i <= n; i++ {
		out <- i
	}
	close(out)
}

func square(in chan int, out chan int) {
	for num := range in {
		out <- num * num
	}
	close(out)
}

type SquaresRequest struct {
	N int `json:"n"`
}

func squaresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req SquaresRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	numbers := make(chan int)
	squares := make(chan int)

	go generate(req.N, numbers)
	go square(numbers, squares)

	result := []int{}
	for sq := range squares {
		result = append(result, sq)
	}

	writeJSON(w, map[string]any{"squares": result})
}

// ==========================================================
// ЗАДАЧА 3. Worker pool: 3 воркера
// ==========================================================

func worker(jobs chan int, results chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range jobs {
		results <- num * 2
	}
}

type DoubleRequest struct {
	Numbers []int `json:"numbers"`
}

func doubleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req DoubleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	jobs := make(chan int)
	results := make(chan int)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	go func() {
		for _, num := range req.Numbers {
			jobs <- num
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	output := []int{}
	for res := range results {
		output = append(output, res)
	}

	writeJSON(w, map[string]any{"doubled": output})
}

// ==========================================================
// ЗАДАЧА 4. Fan-in: три генератора в один канал
// ==========================================================

func gen1(out chan int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

func gen2(out chan int) {
	for i := 10; i <= 50; i += 10 {
		out <- i
	}
	close(out)
}

func gen3(out chan int) {
	for i := 100; i <= 500; i += 100 {
		out <- i
	}
	close(out)
}

func merge(channels ...chan int) chan int {
	merged := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c chan int) {
			defer wg.Done()
			for num := range c {
				merged <- num
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func mergeHandler(w http.ResponseWriter, r *http.Request) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go gen1(ch1)
	go gen2(ch2)
	go gen3(ch3)

	merged := merge(ch1, ch2, ch3)

	result := []int{}
	for num := range merged {
		result = append(result, num)
	}

	writeJSON(w, map[string]any{"merged": result})
}

// ==========================================================
// ЗАДАЧА 5. Семафор на 2 + отмена через context
// ==========================================================

type ProcessRequest struct {
	Numbers []int `json:"numbers"`
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	results := make(chan int, len(req.Numbers))
	sem := make(chan struct{}, 2)

	var wg sync.WaitGroup

	for _, num := range req.Numbers {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return
			}

			if ctx.Err() != nil {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(300 * time.Millisecond):
			}

			res := n * 10

			if res > 50 {
				cancel()
				return
			}

			select {
			case results <- res:
			case <-ctx.Done():
			}
		}(num)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	collected := []int{}
	for res := range results {
		collected = append(collected, res)
	}

	writeJSON(w, map[string]any{
		"collected": collected,
		"canceled":  ctx.Err() != nil,
	})
}

// ==========================================================
// main
// ==========================================================

func main() {
	http.HandleFunc("/sum", sumHandler)
	http.HandleFunc("/squares", squaresHandler)
	http.HandleFunc("/double", doubleHandler)
	http.HandleFunc("/merge", mergeHandler)
	http.HandleFunc("/process", processHandler)

	// раздаём статику из папки static
	http.Handle("/", http.FileServer(http.Dir("./static")))

	println("Server: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}