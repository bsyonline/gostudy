package main

import "fmt"

func main() {
	r := fr(5)
	fmt.Printf("r: %v\n", r)
}

func fr(a int) int {
	if a == 1 {
		return 1
	} else {
		return a * fr(a-1)
	}
}
