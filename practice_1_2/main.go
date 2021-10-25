/*
 * @time 2021/10/22 08:24
 * @version 1.00
 * @author huangsiyi
 */
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("start producer-consumer program")
	ch := make(chan int, 10)
	producer(ch)
	consumer(ch)
	fmt.Println("Enter Ctrl C to quit")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}

func producer(ch chan<- int) {
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
			fmt.Println("data in channel: ", i)
			time.Sleep(time.Millisecond)
		}
		close(ch)
		fmt.Println("close channel")
	}()
}

func consumer(ch <-chan int) {
	go func() {
		for {
			res, ok := <-ch
			if ok {
				fmt.Println("consumer from ch: ", res)
			} else {
				fmt.Println("consumer jam")
				break
			}
			time.Sleep(time.Millisecond)
		}
	}()
}
