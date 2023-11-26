package main

import "fmt"

// 执行顺序：变量赋值>init>main
func main() {
	fmt.Printf("main: abc=%v\n", abc)
}

var abc int = initA()

func init() {
	fmt.Printf("init2\n")
}

func init() {
	fmt.Printf("init1\n")
}

func initA() int {
	fmt.Printf("set abc=100\n")
	return 100
}
