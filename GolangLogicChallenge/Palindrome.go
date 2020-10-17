package main

import "fmt"

func isPalindrome(str string) bool {
	i := 0
	j := len(str) - 1

	for i < j{
		fmt.Println(string(str[i])+" - "+string(string(str[j])))
		if string(str[i]) != string(str[j]) {
			return false
		}
		i++
		j--
	}
	return true
}

func main(){
	input := "madam"
	if isPalindrome(input) == true{
		fmt.Println("Yes")
	} else{
		fmt.Println("No")
	}
}
