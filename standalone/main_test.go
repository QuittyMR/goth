package standalone

import (
	"fmt"
	"gauth"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var jwtService gauth.JWTProvider

func TestMain(m *testing.M) {
	//TODO: Mock DB
	loadConfig()
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/test_auth")
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	defer connection.Close(context.Background())

	jwtService, err = gauth.NewJWTProvider()
	if err != nil {
		panic(err)
	}
	testRecorder = httptest.NewRecorder()
	ginServer := gin.Default()
	authService := gauth.NewAuthService(
		NewPostgresUserService(connection),
		NewPostgresAuthorizationService(connection),
		jwtService,
	)

	ginServer.Use(authService.JWTMiddlware)
	_ = gauth.GetAuthGroup(
		ginServer,
		authService,
	)

	testServer = httptest.NewTLSServer(ginServer)
	testClient = NewClient()
	defer testServer.Close()

	os.Exit(m.Run())
}

func TestAuthService_BasicLogin(t *testing.T) {
	tests := []struct {
		testName string
		username string
		password string
		expected int
	}{
		{"correct", "t4@example.com", "1234", 200},
		{"re-login", "t1@example.com", "1234", 200},
		{"incorrect password", "t1@example.com", "12345", 403},
		{"missing user", "wrong@example.com", "1234", 403},
	}
	t.Parallel()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			response, err := testClient.Post(
				fmt.Sprintf("%s/auth/login/basic", testClient.url),
				"application/json",
				strings.NewReader(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, test.username, test.password)),
			)
			if err != nil {
				t.Fatal("Error getting response from server: ", err.Error())
			}
			assert.EqualValues(t, test.expected, response.StatusCode)
			if response.StatusCode == 200 && test.expected == 200 {
				token, err := testClient.GetToken()
				if err != nil {
					t.Fatal("Error getting data from token:", err.Error())
				}
				assert.Equal(t, test.username, token.Username)
			}
		})
	}

}

//func TestTokenRenewal(t *testing.T) {
//	godotenv.Write()
//}
