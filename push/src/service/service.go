package service

import (
	"fmt"
	"net"
	"net/http"
	"push/proto"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

// TODO: 为了满足拓展性，其实可以用redis存连接的映射关系，即user---serverMachineInstanceID
type server struct {
	proto.UnimplementedPushServiceServer
	connections sync.Map // 用于存储 WebSocket 连接池
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源连接
	},
}

func (s *server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("failed to upgrade websocket: %v\n", err)
		return
	}
	defer conn.Close()

	// 获取用户ID（假设通过查询参数传递）
	userID := r.URL.Query().Get("userName")
	if userID == "" {
		fmt.Println("userID is required")
		return
	}

	// 存储连接
	s.connections.Store(userID, conn)
	fmt.Printf("User %s connected\n", userID)
	conn.WriteJSON([]string{"Welcome to the WebSocket server!"})
	// 监听 WebSocket 消息（可选）
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("User %s disconnected: %v\n", userID, err)
			s.connections.Delete(userID)
			break
		}
	}
}

func StartPushServer() {
	listen, err := net.Listen("tcp", ":8083")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	srv := &server{}
	proto.RegisterPushServiceServer(s, srv)

	defer func() {
		s.Stop()
		listen.Close()
	}()

	// 启动 WebSocket 服务
	//TODO:其实ws可以不放在server里，结构混乱
	go func() {
		http.HandleFunc("/ws", srv.WebSocketHandler)
		http.ListenAndServe(":8084", nil)
		fmt.Println("WebSocket server started on :8084")
	}()

	// 启动 gRPC 服务
	fmt.Println("Serving 8083...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
