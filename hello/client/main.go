package main

import (
	"awesomeProject1/proto/hello"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {

	// TLS连接  记得把server name改成你写的服务器地址
	creds, err := credentials.NewClientTLSFromFile("keys/server.pem", "www.testgrpc.grpc")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	// 连接
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := hello.NewHelloClient(conn)

	// 调用方法
	req := &hello.HelloRequest{Name: "gRPC"}
	res, err := c.SayHello(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}

	fmt.Println(res.Message)
}
