package main

import (
	"fmt"
)

// 延迟执行，用于释放资源，类似栈
func main() {
	defer fmt.Printf("\"step1\": %v\n", "step1")
	defer fmt.Printf("\"step2\": %v\n", "step2")
	fmt.Printf("\"step3\": %v\n", "step3")
	fmt.Printf("\"step4\": %v\n", "step4")
}
