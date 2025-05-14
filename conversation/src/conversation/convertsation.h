#ifndef CONVERSATION_H
#define CONVERSATION_H

#include <sys/socket.h>
#include <unistd.h>
#include "generated/protobuf/conversation.pb.cc"

class ConversationServer {
    public:
    ErrorCode ConversationServer::createSingleConversation(const CreateSingleConversationRequest &req, CreateSingleConversationResponse &resp);

    ErrorCode ConversationServer::createGroupConversation(const CreateGroupConversationRequest, CreateGroupConversationResponse &resp);
    // 处理请求
    bool handleRequest(int socket);
};
    
#endif // CONVERSATION_H