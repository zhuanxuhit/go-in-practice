package main

import (
	"io/ioutil"
	"os"
	"strings"
	"log"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/golang/protobuf/proto"
	"fmt"
)

func Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Print("protoc-gen-go: error:", s)
	os.Exit(1)
}
func Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("protoc-gen-go: error:", s)
	os.Exit(1)
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Error(err, "reading input")
	}
	request := new(plugin.CodeGeneratorRequest)
	if err  := proto.Unmarshal(data, request); err != nil{
		Error(err, "parsing input proto")
	}
	if len(request.FileToGenerate) == 0 {
		Fail("no files to generate")
	}
	//fmt.Println(request)
	//fmt.Println(strings.Join(request.FileToGenerate, ""))
	//fmt.Println("params: " + request.GetParameter())
	// 参数是什么？
	//ouput := make(map[string]string)
	f := request.ProtoFile[0]
	for _, e := range f.EnumType {
		fmt.Println(e)
	}
}
