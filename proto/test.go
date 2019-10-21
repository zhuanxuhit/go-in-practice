package main

import (
	"github.com/golang/protobuf/proto"
	example "github.com/zhuanxuhit/go-in-practice/proto/test"
	"log"
	"fmt"
	"io/ioutil"
)

func main() {
	test1 := &example.Test1{
		A: 150,
	}
	data, err := proto.Marshal(test1)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	fmt.Println(len(data))
	if err := ioutil.WriteFile("/tmp/data1", data, 0644); err != nil {
		log.Fatalln("Failed to write test1:", err)
	}
	//fmt.Println(proto.MessageType("example.Test1"))
	test2 := &example.Test1{}
	proto.Unmarshal(data,test2)
	fmt.Println(test1.A == test2.A)
}
