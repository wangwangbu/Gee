package main

import (
	"context"
	"flag"
	pb "grpc-demo/proto"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

var port string

func init() {
	flag.StringVar(&port, "p", "8000", "启动端口号")
	flag.Parse()
}

type GreeterServer struct{
	pb.UnimplementedGreeterServer
}

// Unary RPC
func (s *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello.world " + r.Name}, nil
}

// Server-side streaming side RPC
func (s *GreeterServer) SayList(r *pb.HelloRequest, stream pb.Greeter_SayListServer) error {
	for i := 0; i<=6; i++ {
		_ = stream.Send(&pb.HelloReply{Message: "hello.list " + strconv.Itoa(i)})
	}
	return nil
}

// Client-side streaming side RPC
func (s *GreeterServer) SayRecord(stream pb.Greeter_SayRecordServer) error {
	reply := "Hello: "
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloReply{Message: "say.record" + reply})
		}
		if err != nil {
			return err
		}
		reply += resp.GetName()
		log.Printf("resp: %v", resp)
	}
	return nil
}

func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "！")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

// Bidirectional streaming RPC
func (s *GreeterServer) SayRoute(stream pb.Greeter_SayRouteServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		reply := magic(in.GetName())

		if err := stream.Send(&pb.HelloReply{Message: reply}); err != nil {
			return err
		}
	}
}

func main() {
	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{})
	lis, _ := net.Listen("tcp", ":"+port)
	server.Serve(lis)
}