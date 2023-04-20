package geecache

import (
	"context"
	"fmt"
	"geecache/consistenthash"
	pb "geecache/geecachepb"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcGetter struct {
	addr string
}

func (g *grpcGetter) Get(in *pb.Request, out *pb.Response) error {
	conn, err := grpc.Dial(g.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewGroupCacheClient(conn)
	resp, err := client.Get(context.Background(), in)
	out.Value = resp.Value
	return err
}

var _ PeerGetter = (*grpcGetter)(nil)

// Server
type GrpcPool struct {
	pb.UnimplementedGroupCacheServer

	self string
	mu sync.Mutex
	peers *consistenthash.Map
	grpcGetters map[string]*grpcGetter
}

func NewGrpcPool(self string) *GrpcPool {
	return &GrpcPool {
		self: self,
		peers: consistenthash.New(defaultReplicas, nil),
		grpcGetters: map[string]*grpcGetter{},
	}
}

func (p *GrpcPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers.Add(peers...)
	for _, peer := range peers {
		p.grpcGetters[peer] = &grpcGetter{
			addr: peer,
		}
	}
}

func (p *GrpcPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		return p.grpcGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*GrpcPool)(nil)

func (p *GrpcPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *GrpcPool) Get(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	p.Log("%s %s", in.Group, in.Key)
	resp := &pb.Response{}
	
	group := GetGroup(in.Group)
	if group == nil {
		p.Log("no such group %v", in.Group)
		return resp, fmt.Errorf("no such group %v", in.Group)
	}
	value, err := group.Get(in.Key)
	if err != nil {
		p.Log("get key %v error %v", in.Key, err)
		return resp, err
	}

	resp.Value = value.ByteSlice()
	return resp, nil
}

func (p *GrpcPool) Run() {
	lis, err := net.Listen("tcp", p.self[7:])
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterGroupCacheServer(server, p)
	
	reflection.Register(server)
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}