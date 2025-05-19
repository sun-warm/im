package message

import (
	"context"
	"fmt"
	"gateway/body"
	"gateway/client"
	"gateway/generated/message"
	"gateway/utils"
)

func SendMessage(ctx context.Context, request body.SendMessageRequest, response *body.SendMessageResponse) error {
	fmt.Println("request:", request.SendMessage.Sender)
	messageTime := utils.GetTime()
	req := message.SendMessageRequest{
		Message: &message.Message{
			Sender:      request.SendMessage.Sender,
			Content:     request.SendMessage.Content,
			MessageTime: messageTime,
		},
		Receiver: request.Receiver,
	}
	resp, err := client.MessageServiceClient.Client.SendMessage(ctx, &req)
	if err != nil {
		return err
	}
	response.StatusCode = int(resp.ErrorCode)
	return nil
}
