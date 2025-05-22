protoc --go_out=. --go-grpc_out=. protobuf/push.proto


# 访问consul
consul agent -dev
# 启动后可以在localhost:8500访问consul