package gauth

import (
	"github.com/gin-gonic/gin"
)

//TODO: Swagger?

func GetAuthGroup(server *gin.Engine, service AuthService) *gin.RouterGroup {
	RolePermissionMap.Refresh(service.authorization)
	authGroup := server.Group("auth")
	authGroup.GET("roles/:userIdentifier", HasPermission("get_roles_self"), service.GetRolesForUser)
	authGroup.GET("roles", HasPermission("get_roles"), service.GetRoles)
	authGroup.GET("permissions/:roleIdentifier", HasPermission("get_permissions"), service.GetPermissionsForRole)
	authGroup.GET("permissions", HasPermission(""), service.GetPermissions)

	loginGroup := authGroup.Group("login")
	loginGroup.POST("basic", service.BasicLogin)

	return authGroup
}

func GetJWTMiddleware(service AuthService) func(ctx *gin.Context) {
	return service.JWTMiddlware
}
