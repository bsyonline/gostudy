package main

import "fmt"

func main() {
	var user User
	user.name = "tom"
	user.age = 20
	fmt.Printf("user: %v\n", user)
	fmt.Printf("user.name: %v\n", user.name)
	fmt.Printf("user.age: %v\n", user.age)

	// 匿名结构体
	var employee struct {
		id   int
		name string
	}

	employee.id = 1
	employee.name = "zhangsan"
	fmt.Printf("employee: %v\n", employee)

	// 结构体初始化
	user1 := User{
		"lisi",
		20,
	}
	fmt.Printf("user1: %v\n", user1)

	// 对部分值初始化
	user2 := User{
		name: "wangwu",
	}
	fmt.Printf("user2: %v\n", user2)
}

// 结构体定义
type User struct {
	name string
	age  int
}
