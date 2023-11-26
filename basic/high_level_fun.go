package main

import "fmt"

func main() {
	r := func1("zhangsan", func2)
	fmt.Printf("%v\n", r)
	f := f3()
	r2 := f("lisi")
	fmt.Printf("%v\n", r2)
}

func func1(name string, f func(string) string) string {
	return func2(name)
}

func func2(name string) string {
	return "hello, " + name
}

func f3() func(string) string {
	return func2
}
