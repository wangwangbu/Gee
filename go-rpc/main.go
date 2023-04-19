package main

import (
	"flag"
	pb "grpc-demo/proto"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}