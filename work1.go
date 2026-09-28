package main

import(
	"fmt"
	"sync"
)

func ParallelSum(nums []int, n int) int{
	if len(nums) == 0{ //обрабатываем  граничные случаи()
		return 0
	}

	if n <= 0{
		n = 1
	}

	if n > len(nums){ //в случае если запросим больше частей, чем элементов. тогда уменьшим n до длины среза
		n = len(nums)
	}

	chunkSize := (len(nums) + n - 1) / n // целочисленное деление с округлением по формуле (a+b-1)/b


	//синхронизация
	var(
		wg	sync.WaitGroup //ожидаем завершение горутин
		mu	sync.Mutex //защищаем общую переменную тотал
		total	int // итог суммы
	)

	for i:= 0; i<n; i++{
		start := i*chunkSize
		end := start+chunkSize
		if end > len(nums){
			end = len(nums)
		}
		if start>=end{
			break
		}

		wg.Add(1) // запускаем горутину
		go func(part []int){ // передача среза как аргумента
			defer wg.Done()
			local := 0
			for _, v := range part{
				local += v
			}
			mu.Lock()
			total += local
			mu.Unlock()
		}(nums[start:end])
	}

	wg.Wait()
	return total
}

func main(){ //тут создаем срез из ляма чисел
	nums := make([]int, 1_000_000)
	for i := range nums{
		nums[i] = i+1
	}

	sum := ParallelSum(nums, 8)
	fmt.Println("Сумма", sum)
}
