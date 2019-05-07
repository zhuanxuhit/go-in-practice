package main

import (
	"errors"
	"fmt"
	"time"
)

func caller() {
	fmt.Println("Enter function caller.")
	panic(errors.New("something wrong")) // 正例。
	panic(fmt.Println)                   // 反例。
	fmt.Println("Exit function caller.")
}

func callerGoPanic() {
	go func() {
		defer func() {
			fmt.Printf("panic: %s\n", recover())
		}()
		panic(errors.New("in caller goroutine"))
	}()
}

func main() {
	// 往panic传递可序序列化的值
	skip := false
	{
		skip = true
		if !skip {
			fmt.Println("Enter function main.")
			caller()
			fmt.Println("Exit function main.")
		}

	}
	// defer 配合 recover
	{
		skip = true
		if !skip {
			fmt.Println("Enter function main.")

			defer func() {
				fmt.Println("Enter defer function.")

				// recover函数的正确用法。
				if p := recover(); p != nil {
					fmt.Printf("panic: %s\n", p)
				}

				fmt.Println("Exit defer function.")
			}()

			// recover函数的错误用法。
			fmt.Printf("no panic: %v\n", recover())

			// 引发panic。
			panic(errors.New("something wrong"))
			// 不会执行到这里了。
			// recover函数的错误用法。
			p := recover()
			fmt.Printf("panic: %s\n", p)

			fmt.Println("Exit function main.")
		}

	}
	// defer是一个fifo
	// 这个都算是一个函数，goroutine 中 panic 没有捕获，也会导致程序出错
	{
		defer fmt.Println("first defer")
		for i := 0; i < 3; i++ {
			defer fmt.Printf("defer in for [%d]\n", i)
		}
		defer fmt.Println("last defer")
	}
	// goroutine 中都需要recover，防止panic
	{
		defer func() {
			fmt.Println("Enter defer2 function.")

			// recover函数的正确用法。
			if p := recover(); p != nil {
				fmt.Printf("panic: %s\n", p)
			}

			fmt.Println("Exit defer2 function.")
		}()

		callerGoPanic()
		select {
		case <-time.After(2 * time.Second):
		}
		println("exist ")
	}
}
