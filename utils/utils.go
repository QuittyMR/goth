package utils

import (
	"github.com/gin-gonic/gin"
	"os"
)

func Alive(context *gin.Context) {
	context.JSON(200, gin.H{"message": "Alive!"})
}

func GetEnvOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

var GetEnv = os.Getenv

type Set map[string]bool

func (set Set) Add(items ...string) {
	for _, item := range items {
		set[item] = true
	}
}

func (set Set) Delete(item string) {
	delete(set, item)
}

func (set Set) Has(item string) bool {
	_, ok := set[item]
	return ok
}

func (set Set) AsList() []string {
	keys := make([]string, 0, len(set))
	for key, _ := range set {
		keys = append(keys, key)
	}
	return keys
}
