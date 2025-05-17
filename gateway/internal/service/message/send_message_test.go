package message

import (
	"context"
	"fmt"
	"gateway/client"
	"gateway/generated/message"
	"gateway/utils"
	"testing"
)

func TestSendMessage(t *testing.T) {
	client.InitMessageClient()
	fmt.Println(1111)
	ctx := context.Background()
	req := message.SendMessageRequest{
		Message: &message.Message{
			Sender:      "user1",
			Content:     "Hello",
			MessageTime: utils.GetTime(),
		},
		Receiver: "user2",
	}
	_, err := client.MessageServiceClient.Client.SendMessage(ctx, &req)
	if err != nil {
		fmt.Println(err)
	}
}
