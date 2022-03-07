package login

import (
	"fmt"
	"gauth/utils"
	"github.com/gin-gonic/gin"
	"github.com/hlandau/passlib"
)

type baseRequestClass struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	storedPassword string
}

func login(credentials baseRequestClass) (newHash string, ok bool) {
	newHash, err := passlib.Verify(credentials.Password, credentials.storedPassword)
	if err != nil {
		return "", false
	}
	return newHash, true
}

func (svc LoginService) BasicLogin(c *gin.Context) {
	request := baseRequestClass{}

	err := c.Bind(&request)
	if err != nil {
		c.JSON(400, &utils.Message{Message: "Missing parameters or malformed request"})
		return
	}

	password, data := svc.user.Get(request.Username)
	request.storedPassword = password

	if _, ok := login(request); ok {
		jwtToken, err := svc.jwt.generateToken(JWTCustomData{Data: data})
		if err != nil {
			c.JSON(500, &utils.Message{Message: fmt.Sprintf("Failed generating JWT token: %s", err.Error())})
			return
		}
		c.SetCookie("session", jwtToken, 60*60*24, "", "", false, true)
		c.JSON(200, &utils.Message{Message: "Authenticated!"})
	} else {
		c.JSON(403, &utils.Message{Message: "Password mismatch!!"})
	}
}

func (svc LoginService) GetRolesForUser(c *gin.Context) {
	//TODO: permission check
	userIdentifier := c.Param("userIdentifier")
	roles := svc.role.GetForUser(userIdentifier)
	c.JSON(200, roles)
}

func (svc LoginService) GetRoles(c *gin.Context) {
	//	TODO: admin check
	roles := svc.role.GetAll()
	c.JSON(200, roles)
}

type LoginService struct {
	user User
	role Roles
	jwt  JWTService
}

func NewLoginService(userService User, roleService Roles, jwt JWTService) LoginService {
	_ = passlib.UseDefaults("20180601") // Argon2I
	return LoginService{userService, roleService, jwt}
}

type User interface {
	Get(identifier string) (password string, userData map[string]interface{})
}

type Roles interface {
	GetForUser(userIdentifier string) []string
	GetAll() []string
}

type Permissions interface {
	GetForRole(roleIdentifier string) []string
	GetAll() []string
}
