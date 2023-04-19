package main

import (
	"context"
	"flag"
	pb "grpc-demo/proto"
	"log"

	"google.golang.org/grpc"
)

var port, addr string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.StringVar(&addr, "addr", "127.0.0.1", "the address to connect to")
	flag.Parse()
}

func main() {
	conn, _ := grpc.Dial(":"+port, grpc.WithInsecure())
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	_ = SayHello(client)
}

func SayHello(client pb.GreeterClient) error {
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "wwb"})
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}