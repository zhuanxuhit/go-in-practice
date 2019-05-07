package main

import (
	"bytes"
	"encoding/binary"
	"github.com/kr/pretty"
	"log"
)

type MyData struct {
	A int64
	B int64
}

//type MyData struct {
//	A int
//	B int
//}
// 注意： 如果字段中有不确定大小的类型，如 int，slice，string 等，则会报错。
// binary write error:binary.Write: invalid type main.MyData

func main() {
	data1 := MyData{
		A: 1,
		B: 2,
	}
	//dataMap := map[int8]string{
	//	1: "MyData",
	//}
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, int8(1)); err != nil {
		log.Fatal("binary write error:", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, data1); err != nil {
		log.Fatal("binary write error:", err)
	}
	pretty.Println(buf)

	var data2 MyData
	var dataNum int8
	if err := binary.Read(buf, binary.LittleEndian, &dataNum); err != nil {
		log.Fatal("binary read error:", err)
	}

	if err := binary.Read(buf, binary.LittleEndian, &data2); err != nil {
		log.Fatal("binary read error:", err)
	}
	pretty.Println(dataNum, data2)
}
