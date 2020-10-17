package main

import (
	"fmt"
	"math/big"
	"time"
)

func fib(n int) *big.Int {
	a := big.NewInt(0)
	b := big.NewInt(1)

	for i := 0; i<=n; i++ {
		a.Add(a, b)
		a,b = b,a
	}
	return a
}

func main() {
	startTime := time.Now()
	for i := 0; i<=10000; i++ {
		fmt.Printf("fib(%d) is \n %d \n ",i, fib(i))
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("time elapsed is %s \n", elapsedTime)
}