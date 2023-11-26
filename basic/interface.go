package main

import "fmt"

func main() {
	var s Shape
	s = Rectangle{width: 4, length: 5}
	fmt.Printf("s.area(): %v\n", s.area())
}

type Shape interface {
	area() float64
}

type Rectangle struct {
	width, length float64
}

func (r Rectangle) area() float64 {
	return r.width * r.length
}
