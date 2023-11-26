package main

import "fmt"

func main() {
	p := Person{
		id:   1,
		name: "tom",
	}
	f1(p)
	fmt.Printf("p: %v\n", p)
	f2(&p)
	fmt.Printf("p: %v\n", p)

}

type Person struct {
	id   int
	name string
}

// 结构体作为函数参数
func f1(person Person) {
	person.id = 20
	person.name = "kate"
	fmt.Printf("person: %v\n", person)
}

// 结构体指针作为函数参数
func f2(person *Person) {
	person.id = 20
	person.name = "kate"
	fmt.Printf("person: %v\n", person)
}
