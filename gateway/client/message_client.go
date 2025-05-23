package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	pb "gateway/generated/message"
	"gateway/utils"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

// FIXME：flag.String 是 Go 标准库 flag 包提供的一个方法，
// 用于定义命令行参数。它的作用是让程序可以通过命令行传递参数，而不是硬编码在代码中。
// 如果命令行没有给这个参数赋值，则默认使用后面的 value
// 用法：./program -message_addr=localhost:9000
// TODO: 后续可以改为用服务发现来获取
// var messageFlags = flag.NewFlagSet("message", flag.ExitOnError)
var (
	messageAddr = flag.String("message_addr", "localhost:8101", "the address to connect to")
	//messageName = flag.String("name", "world", "Name to greet")
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
	serviceName := "message-service" // Consul 中注册的服务名称
	addr, err := utils.GetServiceAddressFromConsul(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service address from Consul: %v", err)
	}
	// 建立 gRPC 连接
	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 建立 gRPC 连接
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
