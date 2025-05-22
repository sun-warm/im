// Package main implements a client for Greeter service.
package client

import (
	"flag"
	"fmt"
	pb "user/generated/push"
	"user/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//const (
//	defaultName = "world"
//)

// 改为用服务发现来获取
//var (
//	addr = flag.String("addr", "localhost:8102", "the address to connect to")
//	name = flag.String("name", defaultName, "Name to greet")
//)

type PushClient struct {
	conn   *grpc.ClientConn
	Client pb.PushServiceClient
}

var PushServiceClient PushClient

// 初始化消息客户端
func InitPushClient() (*PushClient, error) {
	// 解析命令行参数
	flag.Parse()

	serviceName := "push-service" // Consul 中注册的服务名称
	addr, err := utils.GetServiceAddressFromConsul(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service address from Consul: %v", err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	//建立gRPC连接时阻塞，确保在调用RPC之前连接已经建立，避免第一次请求的延迟
	// conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	// if err != nil {
	//     return nil, err
	// }
	// 创建 gRPC 客户端
	client := pb.NewPushServiceClient(conn)
	PushServiceClient = PushClient{conn: conn, Client: client}
	return &PushClient{conn: conn, Client: client}, nil
}

// Close 关闭 gRPC 连接
func (mc *PushClient) Close() {
	if mc.conn != nil {
		mc.conn.Close()
	}
}
