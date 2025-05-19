package service

import (
	"context"
	"conversation/generated/conversation"
	"conversation/utils"
	"errors"
	"fmt"

	"github.com/bytedance/sonic"
)

func (s *Server) CreateGroupConversation(ctx context.Context, req *conversation.CreateGroupConversationRequest) (*conversation.CreateGroupConversationResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.GroupLeader == "" {
		return nil, errors.New("group leader is empty")
	}
	if len(req.Members) == 0 {
		return nil, errors.New("members is empty")
	}
	//FIXME：对于群聊如何获取conversationID？
	//如果需要短且有序的 ID：用 Redis INCR（方法 1）。
	//如果完全去中心化：用 UUID（方法 2）。
	//如果是分布式系统：用 雪花算法（方法 3）。
	//如果强依赖数据库：用 数据库自增 ID（方法 4）。
	users := req.Members
	if !utils.CheckUsersValid(users) {
		return nil, errors.New("users is invalid")
	}
	//目前确认member中必须包含全部成员，leader和admins也需要在member中
	//目前使用方法1生成
	conversationID, err := utils.GenerateGroupConversationID()
	fmt.Println("conversationID:", conversationID)
	if err != nil {
		return nil, err
	}

	err = utils.CreateEmptyZset(conversationID)
	if err != nil {
		return nil, err
	}

	//FIXME:这里需要明确一下，如果创建一个新聊天，但是不发送信息，目前暂定不需要更新最近会话列表，用户之间聊天还是应该以message为确定依据，创建一个conversation但是不发送聊天应当作无意义行为
	//目前直接确定不需要avatar，参照微信等的展示方式
	ConversationInfo := conversation.Conversation{
		ConversationId:   conversationID,
		ConversationType: "group",
		Members:          users,
		Admins:           req.Admins,
		GroupLeader:      req.GroupLeader,
		ConversationName: req.ConversationName,
	}

	jsonString, err := sonic.MarshalString(ConversationInfo)
	fmt.Println("jsonString:", jsonString)
	//TODO：其实推送消息为主，这种更新最近会话列表是可以写在消息队列中消费的，如果有新消息推送到客户端但是客户端发现消息有缺，可以触发补偿拉取
	//FIXME:更新最近会话列表的行为是不是更新自己的就可以，对方又不知道你创建，给他更新没必要。
	err = utils.UpdateUserRecentConversations([]string{req.GroupLeader}, conversationID)
	if err != nil {
		fmt.Println(err.Error())
	}
	resp := conversation.CreateGroupConversationResponse{
		ErrorCode:        conversation.ConversationErrorCode_OK,
		ConversationInfo: &ConversationInfo,
	}
	fmt.Println(resp)
	return &resp, nil
}
