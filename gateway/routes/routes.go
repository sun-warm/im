package routes

import (
	"gateway/internal/controller"
	conversationController "gateway/internal/controller/conversation"
	messageController "gateway/internal/controller/message"
	"gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/login", controller.Login)
	r.POST("/register", controller.Register)
	// r.POST("/ldap_login", func(c *gin.Context) { login.LDAPLogin(c) })
	// r.POST("/check", func(c *gin.Context) { login.CheckLogin(c) })

	//消息相关接口
	message := r.Group("/api/message")
	message.Use(middleware.CheckLoginMiddleWare())
	message.POST("/send", messageController.SendMessage)
	message.GET("/get", messageController.GetMessage)

	//conversation相关接口
	conversation := r.Group("/api/conversation")
	conversation.Use(middleware.CheckLoginMiddleWare())
	conversation.POST("/create/single", conversationController.CreateSingleConversation)
	conversation.POST("/create/group", conversationController.CreateGroupConversation)
	conversation.GET("/get/recent", conversationController.GetRecentConversation)

}
