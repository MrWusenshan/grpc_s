package main

import (
	"awesomeProject1/proto/hello"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"net"
)

const Address = "127.0.0.1:50052"

type helloService struct{}

var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, req *hello.HelloRequest) (res *hello.HelloResponse, err error) {
	res = new(hello.HelloResponse)
	res.Message = fmt.Sprintf("Hello %s.", req.Name)
	return res, nil
}

func main() {
	listener, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to Listen: %v", err)
	}

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("keys/server.pem", "keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	// 实例化grpc Server, 并开启TLS认证
	server := grpc.NewServer(grpc.Creds(creds))

	//server := grpc.NewServer()

	hello.RegisterHelloServer(server, HelloService)
	fmt.Println("Listen on " + Address + " with TLS")
	server.Serve(listener)
}
