package main

import (
	"fmt"
	"strconv"
)

func main() {
	type myFuncType1 func(a int, b int) string
	var f1, f2 myFuncType1
	f1 = add
	r1 := f1(1, 2)
	fmt.Printf("r1: %v\n", r1)
	f2 = max
	r2 := f2(1, 2)
	fmt.Printf("r2: %v\n", r2)

	type myFuncType2 = func(int, int) string
	var f3, f4 myFuncType2
	f3 = add
	r3 := f3(3, 4)
	fmt.Printf("r3: %v\n", r3)
	f4 = max
	r4 := f4(3, 4)
	fmt.Printf("r4: %v\n", r4)

}

func add(a int, b int) string {
	return strconv.Itoa(a) + " + " + strconv.Itoa(b) + " = " + strconv.Itoa(a+b)
}

func max(a int, b int) string {
	max := 0
	if a < b {
		max = b
	} else {
		max = a
	}
	return "max(" + strconv.Itoa(a) + ", " + strconv.Itoa(b) + ") = " + strconv.Itoa(max)
}
