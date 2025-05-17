package service

import (
	"context"
	"conversation/generated/conversation"
	"conversation/utils"
	"errors"
	"fmt"
)

func (s *Server) CreateGroupConversation(ctx context.Context, req *conversation.CreateGroupConversationRequest) (*conversation.CreateGroupConversationResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	//FIXME：对于群聊如何获取conversationID？
	//如果需要短且有序的 ID：用 Redis INCR（方法 1）。
	//如果完全去中心化：用 UUID（方法 2）。
	//如果是分布式系统：用 雪花算法（方法 3）。
	//如果强依赖数据库：用 数据库自增 ID（方法 4）。
	users := req.Members
	if utils.CheckUsersValid(users) {
		return nil, errors.New("users is invalid")
	}
	//FIXME:这里面的Admin和Member以及GroupLeader的关系需要明确一下
	//目前使用方法1生成
	conversationID, err := utils.GenerateGroupConversationID()
	if err != nil {
		return nil, err
	}

	err = utils.CreateEmptyZset(conversationID)
	if err != nil {
		return nil, err
	}

	//TODO:这里需要明确一下，如果创建一个新聊天，但是不发送信息，会不会更新最近会话列表，其实应该更新
	//TODO:需要一个Avatar？
	ConversationInfo := conversation.Conversation{
		ConversationId:   conversationID,
		ConversationType: "group",
		Members:          users,
		Admins:           req.Admins,
		GroupLeader:      req.GroupLeader,
		ConversationName: req.ConversationName,
	}
	err = utils.UpdateUserRecentConversations(users, conversationID)
	if err != nil {
		fmt.Println(err.Error())
	}
	resp := conversation.CreateGroupConversationResponse{
		ErrorCode:        conversation.ConversationErrorCode_OK,
		ConversationInfo: &ConversationInfo,
	}
	return &resp, nil
}
