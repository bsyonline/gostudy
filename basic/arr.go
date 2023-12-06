package main

import "fmt"

func main() {
	var numbers []int

	numbers = append(numbers, 1)
	numbers = append(numbers, 2)
	numbers = append(numbers, 3)

	fmt.Println(numbers)

	var strs []string

	strs = append(strs, "a")
	strs = append(strs, "b")
	strs = append(strs, "c")

	fmt.Println(strs)

}
