package gauth

import (
	"gauth/login"
	"gauth/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}
}

func New(userService *login.UserService) *gin.Engine{
	loadConfig()

	gin.SetMode(utils.GetEnvOrDefault("GIN_MODE", "debug"))
	server := gin.Default()
	_ = server.SetTrustedProxies([]string{utils.GetEnvOrDefault("ALLOWED_SOURCES", "0.0.0.0")})

	setRoutes(server, userService)
	return server
}
