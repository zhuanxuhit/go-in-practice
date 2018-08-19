package main

func test2(x *[3]int) {
	println("test:", x)
	x[1] += 100
}

func main() {
	x := [3]int{1, 2, 3}
	test2(&x)
	println("main:", &x, x[1])
}