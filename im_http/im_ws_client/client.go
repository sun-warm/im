package im_ws_client

import (
	"fmt"
	"google.golang.org/grpc"
	pb "im_http/proto"
	"log"
)

var IMWSClient pb.IMWebsocketClient

func StartIMWSClient() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	IMWSClient = pb.NewIMWebsocketClient(conn)
	fmt.Println("Start WSClient Success")
}
