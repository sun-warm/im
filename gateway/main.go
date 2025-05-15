package main

import (
	"fmt"
	"gateway/client"
	"gateway/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	//1.设置session中间件，参数sessionId，指的是session的名字，也是cookie的名字
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessionId", store))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://172.20.10.2:8081", "http://192.168.71.154:8081"}, // 允许的前端来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "User-Name"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	fmt.Print("start")
	//store, err := rds.NewStore(10, "tcp", "localhost:6379", "", []byte("secret-key"))
	//if err != nil {
	//	panic(err)
	//}
	// 设置session中间件，参数mysession，指的是session的名字，也是cookie的名字
	// store是前面创建的存储引擎，我们可以替换成其他存储引擎
	//r.Use(sessions.Sessions("mysession", store))
	//两种登陆，还需要一个checkLogin

	routes.RegisterRoutes(r)
	r.Run(":8081") // 监听并在 0.0.0.0:8080 上启动服务
}

func main() {
	// redis.ConnectRedis()
	// if redis.ConnectRedis() != nil {
	// 	return
	// }
	// if dao.InitDB() != nil {
	// 	return
	// }
	if _, err := client.InitMessageClient(); err != nil {
		panic(err)
	}
	if _, err := client.InitConversationClient(); err != nil {
		panic(err)
	}
	Start()
}
