package gauth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hlandau/passlib"
	"log"
)

var _ = passlib.UseDefaults("20180601") // Argon2I

type AuthService struct {
	user          UserService
	authorization AuthorizationService
	jwt           JWTProvider
}

type baseRequestClass struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func login(submittedPassword string, userPassword string) (newHash string, ok bool) {
	newHash, err := passlib.Verify(submittedPassword, userPassword)
	if err != nil {
		return "", false
	}
	return newHash, true
}

func (svc AuthService) BasicLogin(c *gin.Context) {
	request := baseRequestClass{}

	err := c.Bind(&request)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "Missing parameters or malformed request"})
	}

	username, password, publicData, err := svc.user.Get(request.Username)
	if err != nil {
		c.AbortWithStatusJSON(403, gin.H{"message": "unauthorized"})
	}

	if _, ok := login(request.Password, password); ok {
		svc.authenticate(c, username, publicData)
		c.JSON(200, gin.H{"message": "Authenticated!"})
	} else {
		c.AbortWithStatusJSON(403, gin.H{"message": "password mismatch"})
	}
}

func (svc AuthService) authenticate(c *gin.Context, username string, publicData map[string]interface{}) {
	roles := svc.authorization.GetRolesForUsers(username)
	jwtToken, err := svc.jwt.generateToken(JWTCustomData{Username: username, Data: publicData, Roles: roles})
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": fmt.Sprintf("Failed generating JWT token: %s", err.Error())})
	}
	c.SetCookie("session", jwtToken, 60*60*24, "", "", false, false)
}

func (svc AuthService) GetRolesForUser(c *gin.Context) {
	userIdentifier := c.Param("userIdentifier")
	roles := svc.authorization.GetRolesForUsers(userIdentifier)
	c.JSON(200, roles)
}

func (svc AuthService) GetRoles(c *gin.Context) {
	roles := svc.authorization.GetAllRoles()
	c.JSON(200, roles)
}

func (svc AuthService) GetPermissions(c *gin.Context) {
	permissions := svc.authorization.GetAllPermissions()
	c.JSON(200, permissions)
}

func (svc AuthService) GetPermissionsForRole(c *gin.Context) {
	roleIdentifier := c.Param("roleIdentifier")
	permissions := svc.authorization.GetRolePermissionMapping(roleIdentifier)
	c.JSON(200, permissions)
}

func (svc AuthService) JWTMiddlware(context *gin.Context) {
	cookie, _ := context.Cookie("session")
	if cookie == "" {
		context.Next()
		return
	}
	token, err := svc.jwt.ReadToken(cookie)
	if err != nil {
		if err.(*jwt.ValidationError).Is(jwt.ErrTokenExpired) {
			log.Print("[INFO] Refreshing token for ", token.Username)
			username, _, publicData, err := svc.user.Get(token.Username)
			if err != nil {
				context.AbortWithStatusJSON(403, gin.H{"reason": "Unauthorized"})
			}
			svc.authenticate(context, username, publicData)
		} else {
			log.Printf("Invalid token received: %s", err.Error())
		}
	} else {
		context.Set(SESSION_DATA_KEY, token.Data)
		context.Set(SESSION_PERMISSIONS_KEY, token.getPermissions())
		context.Set(SESSION_USERNAME_KEY, token.Username)
	}
	context.Next()
}

func NewAuthService(userService UserService, authorizationService AuthorizationService, jwt JWTProvider) AuthService {
	return AuthService{userService, authorizationService, jwt}
}

type UserService interface {
	//TODO: Refactor to use a model? (can't use for JWT or password key leaks through)
	Get(identifier string) (username, password string, publicData map[string]interface{}, err error)
}

type AuthorizationService interface {
	GetRolesForUsers(userIdentifiers ...string) []string
	GetAllRoles() []string
	GetRolePermissionMapping(roleIdentifiers ...string) map[string][]string
	GetAllPermissions() []string
}
