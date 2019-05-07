SRC_DIR = ./
DST_DIR = ./tutorial
gen:
	protoc -I=${SRC_DIR} --go_out=${DST_DIR} ${SRC_DIR}/addressbook.proto

build:
	go build person.go