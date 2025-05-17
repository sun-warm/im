package service

import (
	"context"
	"conversation/generated/conversation"
	"conversation/utils"
	"errors"
	"fmt"
)

func (s *Server) CreateSingleConversation(ctx context.Context, req *conversation.CreateSingleConversationRequest) (*conversation.CreateSingleConversationResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	users := []string{req.Member, req.UserName}
	if utils.CheckUsersValid(users) {
		return nil, errors.New("users is invalid")
	}
	conversationID := utils.GenerateConversationID(req.UserName, req.Member)
	err := utils.CreateEmptyZset(conversationID)
	if err != nil {
		return nil, err
	}
	err = utils.UpdateUserRecentConversations(users, conversationID)
	if err != nil {
		fmt.Println(err.Error())
	}

	ConversationInfo := conversation.Conversation{
		ConversationId:   conversationID,
		ConversationType: "single",
		Members:          []string{req.UserName, req.Member},
	}

	response := conversation.CreateSingleConversationResponse{
		ErrorCode:        conversation.ConversationErrorCode_OK,
		ConversationInfo: &ConversationInfo,
	}

	return &response, nil
}
