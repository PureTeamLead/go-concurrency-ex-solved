package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	reader(doubler(writer(10)))
}

func writer(maxNum int) chan int {
	ch := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i + 1
		}
		close(ch)
	}()

	return ch
}

func doubler(ch chan int) <-chan int {
	var wg sync.WaitGroup
	doubleCh := make(chan int)

	for v := range ch {
		v := v * 2

		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			time.Sleep(500 * time.Millisecond)
			doubleCh <- v
		}(v)

		go func() {
			wg.Wait()
			close(doubleCh)
		}()

	}

	return doubleCh
}

func reader(ch <-chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
