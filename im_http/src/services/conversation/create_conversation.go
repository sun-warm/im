package conversation

import "github.com/gin-gonic/gin"

func CreateConversation(c *gin.Context) {
	c.JSON(200, "create conversation")
}
