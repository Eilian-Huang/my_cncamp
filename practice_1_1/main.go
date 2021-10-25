package main

import (
	"fmt"
)

func main() {
	stringArray := []string{"I", "am", "stupid", "and", "weak"}
	fmt.Printf("stringArray %+v\n", stringArray)
	for index, value := range stringArray {
		if value == "stupid" {
			stringArray[index] = "smart"
		}
		if value == "weak" {
			stringArray[index] = "smart"
		}
	}
	fmt.Printf("stringArray %+v\n", stringArray)
}
