#!/usr/bin/env bash
go build -o ../bin/uppercase
cat words | sort | uniq | ../bin/uppercase