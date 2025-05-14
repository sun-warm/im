package controller

import (
	"gateway/body"
	"gateway/internal/service/message"
	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	//TODO:check auth
	var (
		req  body.SendMessageRequest
		resp body.SendMessageResponse
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, err.Error())
		return
	}
	err := message.SendMessage(c.Request.Context(), req, &resp)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, resp)
}
