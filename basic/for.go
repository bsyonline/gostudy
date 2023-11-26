package main

func main() {
	for3()
}

func for1() {
	for i := 10; i < 20; i++ {
		println(i)
	}
}

func for2() {
	i := 5
	for ; i < 20; i++ {
		println(i)
	}
}

func for3() {
	i := 15
	for i < 20 {
		println(i)
		i++
	}
}

func for4() {
	for {
		println("run")
	}
}
