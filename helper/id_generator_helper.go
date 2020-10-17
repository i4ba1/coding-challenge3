package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenerateId() string {
	rand.Seed(time.Now().Unix())
	//Only lowercase
	charSet := "abcdedfghijklmnopqrstvwxyz"
	var output strings.Builder
	length := 10
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	fmt.Println(output.String())
	return output.String()
}
