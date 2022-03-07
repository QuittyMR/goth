package gauth

import (
	"gauth/login"
	"github.com/gin-gonic/gin"
)

//TODO: Swagger?

func GetAuthGroup(server *gin.Engine, userBlueprint login.User, roleService login.Roles, jwtService login.JWTService) *gin.RouterGroup {
	loginService := login.NewLoginService(userBlueprint, roleService, jwtService)
	authGroup := server.Group("auth")
	loginGroup := authGroup.Group("login")
	loginGroup.POST("basic", loginService.BasicLogin)
	authGroup.GET("roles/:userIdentifier", loginService.GetRolesForUser)
	return authGroup
}
