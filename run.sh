#!/bin/bash

# 设置项目路径
IM_FRONTEND_PATH="/Users/sunwen/Projects/im/im_frontend"
GATEWAY_PATH="/Users/sunwen/Projects/im/gateway"
MESSAGE_PATH="/Users/sunwen/Projects/im/message"
PUSH_PATH="/Users/sunwen/Projects/im/push"

# 启动前端项目
echo "Starting IM Frontend..."
cd "$IM_FRONTEND_PATH"
# npm install
# npm run build
npm run serve & # 后台运行前端服务

# 启动 Gateway 服务
echo "Starting Gateway Service..."
cd "$GATEWAY_PATH"
go build -o gateway main.go
./gateway & # 后台运行 Gateway 服务

# 启动 Message 服务
echo "Starting Message Service..."
cd "$MESSAGE_PATH"
go build -o message main.go
./message & # 后台运行 Message 服务

# 启动 IM WebSocket 服务
#echo "Starting IM WebSocket Service..."
#cd "$PUSH_PATH"
#go build -o push main.go
#./push & # 后台运行 WebSocket 服务

# 等待所有服务启动
echo "All services are starting..."
wait