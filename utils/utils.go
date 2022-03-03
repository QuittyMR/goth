package utils

import (
	"github.com/gin-gonic/gin"
	"os"
)

//TODO: Add DB check, omit from json
type Message struct {
	Message	string `json:"message"`
}

func Alive(context *gin.Context) {
	context.JSON(200, &Message{Message: "Alive!"})
}

func GetEnvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var GetEnv = os.Getenv

func NotImplemented(context *gin.Context) {
	context.JSON(500, &Message{Message: "Not implemented"})
}