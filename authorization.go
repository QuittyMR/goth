package gauth

import (
	"gauth/utils"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

type rolePermissionMapping struct {
	data  map[string][]string
	mutex sync.Mutex
}

var RolePermissionMap rolePermissionMapping = rolePermissionMapping{
	make(map[string][]string),
	sync.Mutex{},
}

func (mapping *rolePermissionMapping) Set(role string, permissions []string) {
	mapping.mutex.Lock()
	defer mapping.mutex.Unlock()
	mapping.data[role] = permissions
}

func (mapping rolePermissionMapping) Get(role string) (permissions []string, ok bool) {
	permissions, ok = mapping.data[role]
	return
}

func (mapping rolePermissionMapping) Refresh(service AuthorizationService, roles ...string) {
	for role, permissions := range service.GetRolePermissionMapping(roles...) {
		RolePermissionMap.Set(role, permissions)
	}
	log.Print("[INFO] Authorization mapping refreshed")
}

func HasPermission(permission string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if permissions, ok := ctx.Get(SESSION_PERMISSIONS_KEY); ok {
			if permissions.(utils.Set).Has(permission) {
				ctx.Next()
				return
			}
		}
		ctx.JSON(403, gin.H{"message": "Unauthorized"})
		return
	}
}

//func (svc AuthService) GetCurrentPermissions(claims jwtStructure) utils.Set {
//	return svc.authorization.GetRolePermissionMapping(claims.Roles...)
//}
