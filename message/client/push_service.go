// Package main implements a client for Greeter service.
package client

import (
	"flag"
	pb "message/generated/push_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

// TODO: 后续可以改为用服务发现来获取
var (
	addr = flag.String("addr", "localhost:8083", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

type PushClient struct {
	conn   *grpc.ClientConn
	Client pb.PushServiceClient
}

var PushServiceClient PushClient

// 初始化消息客户端
func InitMessageClient() (*PushClient, error) {
	// 解析命令行参数
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
