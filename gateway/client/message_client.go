package client

import (
	"context"
	"errors"
	"flag"
	pb "gateway/generated/message"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

// TODO: 后续可以改为用服务发现来获取
var (
	addr = flag.String("addr", "localhost:8101", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

type MessageClient struct {
	conn   *grpc.ClientConn
	Client pb.MessageServiceClient
}

var MessageServiceClient MessageClient

// 初始化消息客户端
func InitMessageClient() (*MessageClient, error) {
	// 解析命令行参数
	flag.Parse()

	// 建立 gRPC 连接
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 建立 gRPC 连接
	conn, err := grpc.DialContext(ctx, *addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	// 验证连接是否可用
	if ok := conn.WaitForStateChange(ctx, conn.GetState()); !ok {
		conn.Close()
		return nil, errors.New("failed to connect to server")
	}

	//建立gRPC连接时阻塞，确保在调用RPC之前连接已经建立，避免第一次请求的延迟
	// conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	// if err != nil {
	//     return nil, err
	// }
	// 创建 gRPC 客户端
	client := pb.NewMessageServiceClient(conn)
	MessageServiceClient = MessageClient{conn: conn, Client: client}
	return &MessageClient{conn: conn, Client: client}, nil
}

// Close 关闭 gRPC 连接
func (mc *MessageClient) Close() {
	if mc.conn != nil {
		mc.conn.Close()
	}
}
