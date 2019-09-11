#!/bin/sh
export GOOS=linux
export GOPACH=amd64
echo 'building...'
#go build -ldflags "-s -w" -v -i
go build -x -o simple_project main.go
echo 'publishing...'
#scp autostatbot query:~/autostatbot/