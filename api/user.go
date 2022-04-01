package api

import (
	"chat/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UserRegister(c *gin.Context) {
	var service service.UserRegisterService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Register()
		c.JSON(200, res)
	} else {
		c.JSON(400, ErrorResponse(err))
		logrus.Info(err)
	}

}
