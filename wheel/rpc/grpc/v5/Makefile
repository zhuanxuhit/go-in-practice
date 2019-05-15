SRC_DIR = ./echo
DST_DIR = ./echo
gen:
	protoc -I=${SRC_DIR} ${SRC_DIR}/echo.proto --go_out=plugins=grpc:${DST_DIR}
clean:
	rm -f greeter_client/client greeter_server/server