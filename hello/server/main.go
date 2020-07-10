package main

import (
	"awesomeProject1/proto/hello"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"net"
)

const Address = "127.0.0.1:50052"

type helloService struct{}

var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, req *hello.HelloRequest) (res *hello.HelloResponse, err error) {
	// 解析metada中的信息并验证

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return nil, grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}

	resp := new(hello.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.\nToken info: appid=%s,appkey=%s", req.Name, appid, appkey)

	return resp, nil
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
