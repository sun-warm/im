package service

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"user/generated/user"

	"google.golang.org/grpc"
)

type server struct {
	user.UnimplementedUserServiceServer
}

func StartUserServer() error {
	//如果使用consul的话
	err := registerServiceWithConsul()
	if err != nil {
		log.Fatalf("failed to register service with Consul: %v", err)
		return err
	}
	listen, err := net.Listen("tcp", ":8104")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server{})

	// 注册 gRPC 健康检查服务
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("push-service", grpc_health_v1.HealthCheckResponse_SERVING)

	defer func() {
		s.Stop()
		listen.Close()
	}()

	fmt.Println("message Serving 8104...")
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
	serviceID := "user-service-1"
	serviceName := "user-service"
	serviceAddress := "localhost"
	servicePort := 8104

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
