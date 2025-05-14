#include "conversation.h"
ErrorCode ConversationServer::createSingleConversation(const CreateSingleConversationRequest &req, CreateSingleConversationResponse &resp){
    // 生成会话ID
    std::string conversationId = std::to_string(time(nullptr));
    resp.set_conversation_id(conversationId);
    return ErrorCode::OK;
}

ErrorCode ConversationServer::createSingleConversation(const CreateSingleConversationRequest &req, CreateSingleConversationResponse &resp){
    // 生成会话ID
    std::string conversationId = std::to_string(time(nullptr));
    resp.set_conversation_id(conversationId);
    return ErrorCode::OK;
}

bool ConversationServer::handleRequest(int socket) {
    // 接收请求消息长度
    uint32_t messageLength;
    if (recv(socket, &messageLength, sizeof(messageLength), 0) <= 0) {
        std::cerr << "Failed to receive message length." << std::endl;
        return false;
    }
    messageLength = ntohl(messageLength);

    // 接收请求消息内容
    std::string serializedMessage(messageLength, '\0');
    if (recv(socket, &serializedMessage[0], messageLength, 0) <= 0) {
        std::cerr << "Failed to receive message." << std::endl;
        return false;
    }

    // 反序列化请求消息
    Request request;
    if (!request.ParseFromString(serializedMessage)) {
        std::cerr << "Failed to parse request." << std::endl;
        return false;
    }

    // 根据请求类型调用相应的处理方法
    if (request.has_create_single_conversation_request()) {
        CreateSingleConversationResponse response;
        if (!createSingleConversation(request.create_single_conversation_request(), response)) {
            response.set_error_code(ErrorCode::UNKNOWN_ERROR);
            //response.set_error_message("Failed to send message.");
        } else {
            response.set_error_code(ErrorCode::OK);
        }
        // 发送响应
        std::string serializedResponse;
        response.SerializeToString(&serializedResponse);
        uint32_t responseLength = htonl(serializedResponse.size());
        send(socket, &responseLength, sizeof(responseLength), 0);
        send(socket, serializedResponse.c_str(), serializedResponse.size(), 0);
    } else if (request.has_create_group_convertsation_request()) {
        CreateGroupConversationResponse response;
        if (!createGroupConversation(request.create_group_convertsation_request(), response)) {
            response.set_error_code(ErrorCode::UNKNOWN_ERROR);
            //response.set_error_message("Failed to send message.");
        } else {
            response.set_error_code(ErrorCode::OK);
        }
        // 发送响应
        std::string serializedResponse;
        response.SerializeToString(&serializedResponse);
        uint32_t responseLength = htonl(serializedResponse.size());
        send(socket, &responseLength, sizeof(responseLength), 0);
        send(socket, serializedResponse.c_str(), serializedResponse.size(), 0);
    }  else {
        std::cerr << "Unknown request type." << std::endl;
        return false;
    }

    return true;
}