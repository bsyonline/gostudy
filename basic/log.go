package main

import (
	"fmt"
	"log"
)

func main() {
	println("start")
	log.Println("print log")
	fmt.Printf("log.Flags(): %v\n", log.Flags())
	log.Panic("panic log")
	log.Fatal("has same error")
	println("end")
}
