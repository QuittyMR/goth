package standalone

import (
	"gauth"
	"gauth/utils"
	"github.com/gin-contrib/cors"
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

func getCorsMiddleware() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	if utils.GetEnvOrDefault("GIN_MODE", "debug") == "debug" {
		log.Println("[WARNING] GIN_MODE set to debug.")
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = []string{"http://localhost"}
	}

	return cors.New(corsConfig)
}

func main() {
	loadConfig()
	gin.SetMode(utils.GetEnvOrDefault("GIN_MODE", "debug"))
	server := gin.Default()
	_ = server.SetTrustedProxies([]string{utils.GetEnvOrDefault("ALLOWED_SOURCES", "localhost")})
	server.Use(getCorsMiddleware())
	//TODO: Make JWTProvider instantiation internal
	server.GET("alive", utils.Alive)

	connection, err := pgx.Connect(context.Background(), utils.GetEnv("PSQL_CONNECTION_STRING"))
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	defer connection.Close(context.Background())

	jwtService, _ := gauth.NewJWTProvider()
	authService := gauth.NewAuthService(
		NewPostgresUserService(connection),
		NewPostgresAuthorizationService(connection),
		jwtService,
	)
	_ = gauth.GetAuthGroup(
		server,
		authService,
	)
	server.Use(gauth.GetJWTMiddleware(authService))

	err = server.Run()
	if err != nil {
		log.Fatal("Could not start server:", err.Error())
	}
}
