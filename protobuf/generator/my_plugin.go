package main

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

type JsonFormat struct {
	Typ      string           `json:"type"`
	Filename string           `json:"filename"`
	Name     string           `json:"name"`
	Values   map[string]int32 `json:"values"`
}

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Error(err, "reading input")
	}
	request := new(plugin.CodeGeneratorRequest)
	if err := proto.Unmarshal(data, request); err != nil {
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
	e := f.GetEnumType()[0]

	output := new(JsonFormat)
	output.Name = e.GetName()
	output.Filename = f.GetName()
	output.Typ = "Enum"
	values := make(map[string]int32)
	for _, v := range e.GetValue() {
		values[v.GetName()] = v.GetNumber()
	}
	output.Values = values

	response := new(plugin.CodeGeneratorResponse)

	content, err := json.Marshal(output)
	if err != nil {
		Error(err, "failed to marshal output json")
	}
	//ioutil.WriteFile("test.txt", content, os.ModePerm)
	//fmt.Println(content)
	//os.Exit(1)
	response.File = append(response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String("hello.json"),
		Content: proto.String(string(content)),
	})

	// Send back the results.
	data, err = proto.Marshal(response)
	if err != nil {
		Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		Error(err, "failed to write output proto")
	}
}
