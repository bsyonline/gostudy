package main

import "fmt"

func main() {
	sum := func(a int, b int) int {
		return a + b
	}
	r := sum(1, 2)
	fmt.Printf("r: %v\n", r)

	add := func(a int, b int) int {
		return a + b
	}(3, 4)
	fmt.Printf("add: %v\n", add)
}
