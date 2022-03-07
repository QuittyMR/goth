package main

import (
	"gauth"
	"gauth/internal"
	"gauth/login"
	"gauth/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"log"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}
}

func main() {
	loadConfig()
	connection, err := pgx.Connect(context.Background(), utils.GetEnv("PSQL_CONNECTION_STRING"))
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	defer connection.Close(context.Background())
	userService := internal.NewUserService(connection)
	gin.SetMode(utils.GetEnvOrDefault("GIN_MODE", "debug"))
	server := gin.Default()
	_ = server.SetTrustedProxies([]string{utils.GetEnvOrDefault("ALLOWED_SOURCES", "localhost")})
	server.GET("alive", utils.Alive)
	//TODO: Make JWTService instantiation internal
	jwtService, _ := login.NewJWTService()
	server.Use(jwtService.JWTMiddlware)
	roleService := internal.NewPostgresRoleService(connection)
	_ = gauth.GetAuthGroup(server, userService, roleService, jwtService)

	err = server.Run()
	if err != nil {
		log.Fatal("Could not start server:", err.Error())
	}

}
