package message

import (
	"context"
	"fmt"
	"gateway/client"
	"gateway/generated/message"
	"testing"
)

func TestSendMessage(t *testing.T) {
	client.InitMessageClient()
	fmt.Println(1111)
	ctx := context.Background()
	req := message.SendMessageRequest{Content: "a", UserName: "A", Receiver: "a"}
	_, err := client.MessageServiceClient.Client.SendMessage(ctx, &req)
	if err != nil {
		fmt.Println(err)
	}
}
