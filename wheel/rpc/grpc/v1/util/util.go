package util

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	pb "github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v1/tutorial"
	"io"
)

const HeaderSize = 2

func ReadTCP(rr *bufio.Reader) (message *pb.Message, err error) {
	buf := make([]byte, HeaderSize)
	//n, err := rr.Read(buf)
	n, err := io.ReadFull(rr, buf)
	if err != nil {
		log.Errorf("Read err:", err)
		return
	}
	if n != HeaderSize {
		err = fmt.Errorf("reads only %d bytes", n)
		return
	}
	bodyLen := binary.BigEndian.Uint16(buf)
	log.Infof("will read %d body len", bodyLen)
	// 2. 读取Message
	protoBuf := make([]byte, bodyLen)
	n, err = io.ReadFull(rr, protoBuf)

	if err != nil {
		log.Errorf("Read err:", err)
		return
	}
	if n != int(bodyLen) {
		err = fmt.Errorf("reads only %d bytes", n)
		return
	}
	message = &pb.Message{}
	if err = proto.Unmarshal(protoBuf, message); err != nil {
		log.Errorf("Failed to Unmarshal message", err)
		return
	}
	return
}
func WriteTCP(wr *bufio.Writer, message *pb.Message) (err error) {
	var (
		buf []byte
	)

	out, err := proto.Marshal(message)
	if err != nil {
		log.Errorf("Failed to Marshal message:", err)
		return
	}
	bodyLen := uint16(len(out))
	packLen := HeaderSize + bodyLen
	log.Infof("send total len:%d, head:%d body:%d", packLen, HeaderSize, bodyLen)
	//fmt.Printf("send total len:%d, head:%d body:%d",packLen, HeaderSize, bodyLen)
	log.Flush()

	buf = make([]byte, packLen)
	binary.BigEndian.PutUint16(buf, bodyLen)
	for i := range out {
		buf[i+HeaderSize] = out[i]
	}
	nn, err := wr.Write(buf)
	if err != nil {
		log.Errorf("Write err:", err)
		return
	}
	if nn != int(packLen) {
		err = fmt.Errorf("write only %d bytes", nn)
		return
	}
	wr.Flush()
	return
}
