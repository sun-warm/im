package message

import (
	"context"
	"fmt"
	"gateway/body"
	"gateway/client"
	"gateway/proto"
)

func SendMessage(ctx context.Context, request body.SendMessageRequest, response *body.SendMessageResponse) error {
	fmt.Println("request:", request.UserName)
	resp, err := client.MessageServiceClient.Client.SendMessage(ctx, &proto.SendMessageRequest{})
	if err != nil {
		return err
	}
	response.StatusCode = int(resp.ErrorCode)
	return nil
}
