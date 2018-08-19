package main

func main() {
	var x []int
	for i := 0; i < 1024; i++ {
		x = append(x, i)
	}
}