package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CheckLoginMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if userName, ok := session.Get("user_name").(string); ok && c.GetHeader("User-Name") == userName {
			c.Next()
		} else {
			//s := fmt.Sprintf("userName is %s, session.get is %s", c.GetHeader("User-Name"), session.Get("user_name").(string))
			c.JSON(401, "user has not login")
			//c.JSON(401, s)
			c.Abort()
		}
	}
}
