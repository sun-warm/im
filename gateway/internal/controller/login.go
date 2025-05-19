package controller

import (
	"gateway/body"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	//TODO:check auth
	var req body.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(404, err.Error())
	}
	session := sessions.Default(c)
	session.Set("user_name", req.UserName)
	if err := session.Save(); err != nil {
		c.JSON(500, err.Error())
		return
	}
	c.JSON(200, "OK")
}
