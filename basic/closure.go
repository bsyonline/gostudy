package main

import "fmt"

func main() {
	f := say("hello")
	r := f("zhangsan")
	fmt.Printf("r: %v\n", r)
	f1 := say("goodbye")
	r = f1("lisi")
	fmt.Printf("r: %v\n", r)

	add, sub := cal(100)
	i := add(100) // i=200
	i = add(100)  // i=300
	fmt.Printf("i: %v\n", i)
	i = sub(50) // i=250
	fmt.Printf("i: %v\n", i)
}

func say(op string) func(string) string {
	return func(name string) string {
		return op + ", " + name
	}
}

func cal(base int) (func(int) int, func(int) int) {
	add := func(a int) int {
		base += a
		return base
	}

	sub := func(a int) int {
		base -= a
		return base
	}
	return add, sub
}
