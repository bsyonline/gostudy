package main

import "fmt"

func main() {
	cat := Cat{
		name:  "tom",
		color: "blue",
	}
	cat.run()
}

type Cat struct {
	name  string
	color string
}

// 方法：将函数绑定到结构体
func (c Cat) run() {
	fmt.Printf("%v cat named %v run\n", c.color, c.name)
}
