package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	pb "grpc-demo/proto"
	"io"
	"log"
	"os"
	"strings"
	"time"

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
	// _ = SayHello(client)
	// _ = SayList(client, &pb.HelloRequest{Name: "wwb"})
	// _ = SayRecord(client)
	SayRoute(client)
}

// Unary RPC
func SayHello(client pb.GreeterClient) error {
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{Name: "wwb"})
	log.Printf("client.SayHello resp: %s", resp.Message)
	return nil
}

// Server-side streaming RPC
func SayList(client pb.GreeterClient, r *pb.HelloRequest) error {
	stream, _ := client.SayList(context.Background(), r)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: %v", resp)
	}
	return nil
}

// Client-side streaming RPC
func SayRecord(client pb.GreeterClient) error {
	stream, _ := client.SayRecord(context.Background())
	names := []string{"wwb", "wwb2", "wwb3"}
	for _, name := range names {
		err := stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			log.Fatalf("c.LotsOfGreetings stream.Send(%v) failed, err: %v", name, err)
		}
	}
	resp, _ := stream.CloseAndRecv()
	log.Printf("resp err: %v", resp)
	return nil
}

// Bidirectional streaming RPC
func SayRoute(client pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	stream, err := client.SayRoute(ctx)
	if err != nil {
		log.Fatalf("client.SayRoute failed, err: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.SayRoute stream.Recv() failed, err: %v", err)
			}
			fmt.Printf("AI: %s\n", in.GetMessage())
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			continue
		}
		if strings.ToUpper(cmd) == "QUIT" {
			break
		}
		if err := stream.Send(&pb.HelloRequest{Name: cmd}); err!=nil {
			log.Fatalf("client.SayRoute stream.Send(%v) failed: %v", cmd, err)
		}
	}
	stream.CloseSend()
	<-waitc
}