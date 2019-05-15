package main

import "context"

func main() {
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//a := r.Perm(200)[:5]
	//fmt.Println(a)
	ctx := context.Background()
	<-ctx.Done()
}
