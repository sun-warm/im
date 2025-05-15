package message

import (
	"context"
	"fmt"
	"gateway/body"
	"gateway/client"
	"gateway/generated/message"
)

func SendMessage(ctx context.Context, request body.SendMessageRequest, response *body.SendMessageResponse) error {
	fmt.Println("request:", request.UserName)
	req := message.SendMessageRequest{
		UserName: request.UserName,
		Receiver: request.Receiver,
		Content:  request.Content,
	}
	resp, err := client.MessageServiceClient.Client.SendMessage(ctx, &req)
	if err != nil {
		return err
	}
	fmt.Println(1111111)
	response.StatusCode = int(resp.ErrorCode)
	return nil
}
