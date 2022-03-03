package login

import (
	"fmt"
	"gauth/utils"
	"github.com/gin-gonic/gin"
)

type baseRequestClass struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserService interface {
	GetPassword(email string) string
	SetPassword(email, newPassword string) error
	SetActive(email string, isActive bool) error
}

func (svc LoginService) login(email, password string) bool {
	return svc.user.GetPassword(email) == password
}

func (svc LoginService) BasicLogin(context *gin.Context) {
	//TODO: return JWT
	request := struct {
		baseRequestClass
	}{}

	err := context.Bind(&request)
	if err != nil {
		context.JSON(400, &utils.Message{Message: "Missing parameters or malformed request"})
		return
	}

	if svc.login(request.Email, request.Password) {
		context.JSON(200, &utils.Message{Message: "Authenticated!"})
		return
	}
	context.JSON(403, &utils.Message{Message: "Password mismatch!!"})
}

func (svc LoginService) SetPassword(context *gin.Context) {
	request := struct {
		baseRequestClass
		NewPassword string `json:"new_password" binding:"required"`
	}{}

	err := context.Bind(&request)
	if err != nil {
		context.JSON(400, &utils.Message{Message: fmt.Sprintf("Missing parameters or malformed request: %s", err.Error())})
		return
	}

	if svc.login(request.Email, request.Password) {
		err := svc.user.SetPassword(request.Email, request.NewPassword)
		if err != nil {
			context.JSON(400, &utils.Message{fmt.Sprintf("error attempting to set new password: %s", err.Error())})
		}
	} else {
		context.JSON(403, &utils.Message{"unauthorized"})
	}
}

func (svc LoginService) SetActive(context *gin.Context) {
	request := struct {
		baseRequestClass
		IsActive string `json:"is_active" binding:"required"`
	}{}

	err := context.Bind(&request)
	if err != nil {
		context.JSON(400, &utils.Message{Message: "Missing parameters or malformed request"})
		return
	}

	if svc.login(request.Email, request.IsActive) {
		err := svc.user.SetActive(request.Email, request.IsActive == "t")
		if err != nil {
			context.JSON(400, &utils.Message{fmt.Sprintf("error attempting to set active state: %s", err.Error())})
		}
	} else {
		context.JSON(403, &utils.Message{"unauthorized"})
	}
}

type LoginService struct {
	user UserService
}

func NewLoginService(service *UserService) LoginService {
	return LoginService{*service}
}
