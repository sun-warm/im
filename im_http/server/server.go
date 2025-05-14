package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	middleWare "im_http/src/middleware"
	"im_http/src/services/conversation"
	"im_http/src/services/message"
	"net/http"
	"time"
)

func StartIMHTTPServer() {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionId", store))

	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	r.Use(sessions.Sessions("mysession", store))
	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:8081", "http://172.20.10.2:8081", "http://192.168.71.154:8081"}, // 允许的前端来源
		//AllowOrigins:     []string{"*"},//但是AllowOrigins设置为*时，AllowCredentials就不能设置为true，否则会报错，因此修改为下面的Func
		AllowOriginFunc: func(origin string) bool {
			return true // 允许所有来源
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "User-Name"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, HTTP/2!")
	})
	conversationModule := r.Group("conversation").Use(middleWare.CheckLoginMiddleWare())
	conversationModule.POST("create", func(c *gin.Context) {
		conversation.CreateConversation(c)
	})
	messageModule := r.Group("message").Use(middleWare.CheckLoginMiddleWare())
	messageModule.POST("send", func(c *gin.Context) {
		message.SendMessage(c)
	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务

}
