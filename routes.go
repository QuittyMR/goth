package gauth

import (
	"gauth/login"
	"gauth/utils"
	"github.com/gin-gonic/gin"
)

//TODO: Swagger?

func setRoutes(server *gin.Engine, userService *login.UserService) {
	loginService := login.NewLoginService(userService)
	server.GET("alive", utils.Alive)
	server.POST("login/basic", loginService.BasicLogin)
	server.POST("user/password", loginService.SetPassword)
	server.POST("user/active", loginService.SetActive)
}
