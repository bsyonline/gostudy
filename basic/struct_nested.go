package main

import "fmt"

func main() {
	var boy Boy
	var dog Dog
	dog.name = "gaofei"
	dog.color = "black"
	boy.name = "zhangsan"
	boy.age = 20
	boy.dog = dog
	fmt.Printf("person: %v\n", boy)
}

type Boy struct {
	name string
	age  int
	dog  Dog
}

type Dog struct {
	name  string
	color string
}
