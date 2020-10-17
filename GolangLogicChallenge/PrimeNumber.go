package main

import (
	"fmt"
	"math"
	"math/big"
)

func generatePrimeNumber(n int){
	z := big.NewInt(int64(n))
	for i := 1; i<n; i++ {
		if z.ProbablyPrime(0) {
			fmt.Printf("%d => ", z)
		}
	}
}

func sievePrimeFactors(){
	const n = 10000
	var myArray[n]bool
	for i := 0; i< len(myArray); i++ {
		myArray[i] = true
	}

	for j := 2; j<int(math.Sqrt(n)); j++ {
		if myArray[j] == true{
			for k := j*j; k<n; k += j{
				myArray[k] = false;
			}
		}
	}

	fmt.Println("List of prime numbers upto given number are : ")
	for i := 2; i< len(myArray); i++ {
		if myArray[i] == true {
			fmt.Printf("%d \n =>",i)
		}
	}
}

func main(){
	sievePrimeFactors()
}