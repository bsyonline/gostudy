package main

import "fmt"

func main() {
	fmt.Printf("--------普通指针--------\n")
	var name string
	name = "zhangsan"
	var pname *string //string指针类型
	pname = &name     //指针赋值

	fmt.Printf("name: %v\n", name)
	fmt.Printf("pname: %v\n", pname)
	fmt.Printf("pname指针的值: %v\n", *pname) //指针取值
	fmt.Printf("--------结构体指针指针--------\n")
	type User struct {
		id   int
		name string
		age  int
	}

	var puser *User //结构体指针

	user := User{
		1, "zhansan", 40,
	}
	puser = &user

	fmt.Printf("user: %v\n", user)
	fmt.Printf("puser: %p\n", puser)
	fmt.Printf("puser指针的值: %v\n", *puser)

	// 用new声明结构体指针
	puser1 := new(User)
	(*puser1).id = 3 //*号可以省略
	puser1.name = "jack"
	puser1.age = 30
	fmt.Printf("puser1: %p\n", puser1)
	fmt.Printf("puser1指针的值: %v\n", *puser1)
}
