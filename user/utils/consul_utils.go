package utils

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

// 从 Consul 获取服务地址
func GetServiceAddressFromConsul(serviceName string) (string, error) {
	// 创建 Consul 客户端
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create Consul client: %v", err)
	}

	// 查询服务
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to query Consul for service %s: %v", serviceName, err)
	}

	// 如果没有找到服务实例
	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	// 返回第一个健康的服务实例地址
	service := services[0]
	address := fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	return address, nil
}
