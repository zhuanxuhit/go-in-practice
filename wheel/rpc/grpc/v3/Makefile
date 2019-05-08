SRC_DIR = ./
DST_DIR = ./tutorial
gen:
	protoc -I=${SRC_DIR} --go_out=${DST_DIR} ${SRC_DIR}/message.proto

build:
	go build client.go
	go build server.go
clean:
	rm -f client server