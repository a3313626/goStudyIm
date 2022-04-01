package router

import (
	"chat/api"
	"chat/service"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Recovery(), gin.Logger())
	//Recovery 防止中断程序
	//logger 日志

	v1 := r.Group("/")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "hello")
		})
		v1.POST("user/register", api.UserRegister)
		v1.GET("ws", service.Handler)
	}

	return r
}
