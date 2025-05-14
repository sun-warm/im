package routes

import (
	"gateway/internal/controller"
	"gateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.POST("/login", controller.Login)
	// r.POST("/ldap_login", func(c *gin.Context) { login.LDAPLogin(c) })
	// r.POST("/check", func(c *gin.Context) { login.CheckLogin(c) })

	//报表相关接口
	message := r.Group("/api/message")
	message.Use(middleware.CheckLoginMiddleWare())
	message.POST("/send", controller.SendMessage)

}
