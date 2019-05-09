SRC_DIR = ./helloworld
DST_DIR = ./helloworld
gen:
	protoc -I=${SRC_DIR} ${SRC_DIR}/helloworld.proto --go_out=plugins=grpc:${DST_DIR}

build:
	go build -o greeter_client/client greeter_client/main.go
	go build -o greeter_server/server greeter_server/main.go
clean:
	rm -f greeter_client/client greeter_server/server