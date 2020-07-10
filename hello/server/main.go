package main

import (
	"awesomeProject1/proto/hello"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

	server := grpc.NewServer()

	hello.RegisterHelloServer(server, HelloService)
	//grpclog.Println("Listen on" + Address)
	server.Serve(listener)
}
