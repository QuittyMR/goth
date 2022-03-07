package gauth

//
//import (
//	"gauth/login"
//	"gauth/providers"
//	"gauth/utils"
//	"github.com/gin-gonic/gin"
//	"github.com/joho/godotenv"
//	"log"
//)
//
//func loadConfig() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatalf("Error loading .env file: %s", err.Error())
//	}
//}
//
//func NewServer(userService *providers.UserService) *gin.Engine {
//	loadConfig()
//	gin.SetMode(utils.GetEnvOrDefault("GIN_MODE", "debug"))
//	server := gin.Default()
//	authGroup := GetAuthGroup(server, userService, jwtService)
//	_ = server.SetTrustedProxies([]string{utils.GetEnvOrDefault("ALLOWED_SOURCES", "localhost")})
//	jwtService, _ := login.NewJWTService()
//	server.Use(jwtService.JWTMiddlware)
//
//	return server
//}
