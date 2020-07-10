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
	resp := new(hello.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.", req.Name)

	return resp, nil
}

func main() {
	listener, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to Listen: %v", err)
	}

	var opts []grpc.ServerOption
	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("keys/server.pem", "keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	opts = append(opts, grpc.Creds(creds))

	// 注册interceptor
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	// 实例化grpc Server, 并开启TLS认证
	server := grpc.NewServer(opts...)

	//server := grpc.NewServer()

	hello.RegisterHelloServer(server, HelloService)
	fmt.Println("Listen on " + Address + " with TLS")
	server.Serve(listener)
}

// auth 验证Token
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
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
		return grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}

	return nil
}

// interceptor 拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}
