#!/bin/sh

protoc --plugin=protoc-gen-custom=my_plugin --custom_out=./hello hello.proto
