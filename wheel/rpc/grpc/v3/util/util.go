package util

import (
	"bufio"
	"fmt"
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	pb "github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v3/tutorial"
	"io"
	"reflect"
)

var messageMap = map[pb.Message]reflect.Type{

}

func init() {
	messageMap[pb.Message_Message_XX1Request] = reflect.TypeOf(pb.XX1Request{})
	messageMap[pb.Message_Message_XX2Request] = reflect.TypeOf(pb.XX2Request{})
	messageMap[pb.Message_Message_XX1Response] = reflect.TypeOf(pb.XX1Response{})
	messageMap[pb.Message_Message_XX2Response] = reflect.TypeOf(pb.XX2Response{})
}

func ReadTCP(rr *bufio.Reader, headerSize int) (request interface{}, err error) {
	buf := make([]byte, headerSize)
	//n, err := rr.Read(buf)
	n, err := io.ReadFull(rr, buf)
	if err != nil {
		log.Errorf("Read err:", err)
		return
	}
	if n != headerSize {
		err = fmt.Errorf("reads only %d bytes", n)
		return
	}
	messageHeader := &pb.MessageHeader{
	}
	if err = proto.Unmarshal(buf, messageHeader); err != nil {
		log.Errorf("Failed to Unmarshal message", err)
		return
	}

	log.Infof("read header:%v", messageHeader)

	req := reflect.New(messageMap[messageHeader.MessageType])

	log.Infof("new msg:%+v", req.Type())

	// 下一步就是我们要去读取具体的数据。
	buf = make([]byte, messageHeader.MessageLength)
	n, err = io.ReadFull(rr, buf)
	if err != nil {
		log.Errorf("Read err:", err)
		return
	}
	if int32(n) != messageHeader.MessageLength {
		err = fmt.Errorf("reads only %d bytes", n)
		return
	}

	if _, ok := req.Interface().(proto.Message); !ok {
		err = fmt.Errorf("request:%+v can't convert to proto.Message", request)
		return
	}

	if err = proto.Unmarshal(buf, req.Interface().(proto.Message)); err != nil {
		log.Errorf("Failed to Unmarshal message", err)
		return
	}
	request = req.Interface()
	log.Infof("read msg body:%+v", request)
	return
}
func WriteTCP(wr *bufio.Writer, request proto.Message, messageType pb.Message) (err error) {
	var (
		messageHeader pb.MessageHeader
	)

	outBody, err := proto.Marshal(request)
	if err != nil {
		log.Errorf("Failed to Marshal request:", err)
		return
	}
	messageHeader.MessageType = messageType
	messageHeader.MessageLength = int32(len(outBody))

	outHeader, err := proto.Marshal(&messageHeader)
	if err != nil {
		log.Errorf("Failed to Marshal messageHeader:", err)
		return
	}

	packLen := len(outHeader) + len(outBody)
	log.Infof("send total len:%d, head:%d body:%d", packLen, len(outHeader), len(outBody))

	nn, err := wr.Write(outHeader)
	if err != nil {
		log.Errorf("Write err:", err)
		return
	}
	if nn != int(len(outHeader)) {
		err = fmt.Errorf("write only %d bytes", nn)
		return
	}

	nn, err = wr.Write(outBody)
	if err != nil {
		log.Errorf("Write err:", err)
		return
	}
	if nn != int(len(outBody)) {
		err = fmt.Errorf("write only %d bytes", nn)
		return
	}
	wr.Flush()
	return
}
