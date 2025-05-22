package service

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"push/proto"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

func StartPushServer() error {
	//如果使用consul的话
	err := registerServiceWithConsul()
	if err != nil {
		log.Fatalf("failed to register service with Consul: %v", err)
		return err
	}
	listen, err := net.Listen("tcp", ":8102")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	srv := &server{}
	proto.RegisterPushServiceServer(s, srv)

	// 注册 gRPC 健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("push-service", grpc_health_v1.HealthCheckResponse_SERVING)

	defer func() {
		s.Stop()
		listen.Close()
	}()

	// 启动 WebSocket 服务
	//TODO:其实ws可以不放在server里，结构混乱
	go func() {
		http.HandleFunc("/ws", srv.WebSocketHandler)
		http.ListenAndServe(":8105", nil)
		fmt.Println("WebSocket server started on :8105")
	}()

	// 启动 gRPC 服务
	fmt.Println("Serving 8083...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return err
	}
	return nil
}

func registerServiceWithConsul() error {
	// 创建 Consul 客户端
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return fmt.Errorf("failed to create Consul client: %v", err)
	}

	// 服务注册信息
	serviceID := "push-service-1"
	serviceName := "push-service"
	serviceAddress := "localhost"
	servicePort := 8102

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: serviceAddress,
		Port:    servicePort,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", serviceAddress, servicePort),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	// 注册服务到 Consul
	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service with Consul: %v", err)
	}

	log.Printf("Service %s registered with Consul", serviceName)
	return nil
}
