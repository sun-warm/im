# 生成proto文件
protoc --go_out=./generated --go-grpc_out=./generated protobuf/message/message.proto

# 生成proto文件
protoc --go_out=./generated --go-grpc_out=./generated protobuf/conversation/conversation.proto

# 生成proto文件
protoc --go_out=./generated --go-grpc_out=./generated protobuf/user/user.proto