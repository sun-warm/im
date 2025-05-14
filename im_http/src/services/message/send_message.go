package message

import (
	"context"
	"github.com/gin-gonic/gin"
	"im_http/im_ws_client"
	pb "im_http/proto"
)

func SendMessage(c *gin.Context) {
	ctx := context.Background()
	resp, err := im_ws_client.IMWSClient.ServerPushMsgToUser(ctx, &pb.ServerPushMsgToUserRequest{
		UserID: "",
		MessageContent: &pb.MessageContent{
			Message: "",
		},
	})
	if err != nil {
		c.JSON(500, "send message failed")
		return
	}
	if resp.StatusCode == 200 {
		c.JSON(200, "send message success")
	} else {
		c.JSON(500, "send message failed")
	}
}
