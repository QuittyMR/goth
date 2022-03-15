package standalone

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"goth"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var jwtService goth.JWTProvider

func TestMain(m *testing.M) {
	//TODO: Mock DB
	loadConfig()
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/test_auth")
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	defer connection.Close(context.Background())

	jwtService, err = goth.NewJWTProvider(true)
	if err != nil {
		panic(err)
	}
	testRecorder = httptest.NewRecorder()
	ginServer := gin.Default()
	authService := goth.NewAuthService(
		NewPostgresUserService(connection),
		NewPostgresAuthorizationService(connection),
		jwtService,
	)

	ginServer.Use(authService.JWTMiddlware)
	_ = goth.GetAuthGroup(
		ginServer,
		authService,
	)

	testServer = httptest.NewTLSServer(ginServer)
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
	testClient := NewClient(t)
	t.Parallel()
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			response, err := testClient.Post(
				testClient.Url("login/basic"),
				"application/json",
				strings.NewReader(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, test.username, test.password)),
			)
			if err != nil {
				t.Fatal("Error getting response from server: ", err.Error())
			}
			assert.EqualValues(t, test.expected, response.StatusCode)
			if response.StatusCode == 200 && test.expected == 200 {
				token := testClient.GetToken()
				assert.Equal(t, test.username, token.Username, "Username should appear in token")
			}
		})
	}

}

func TestTokenRenewal(t *testing.T) {
	testClient := NewClient(t)
	response, err := testClient.Post(
		testClient.Url("login/basic"),
		"application/json",
		strings.NewReader(`{"username":"t1@example.com", "password":"1234"}`),
	)
	if err != nil {
		t.Fatal("Error getting response from server: ", err.Error())
	}
	assert.EqualValues(t, 200, response.StatusCode)
	originalToken := response.Cookies()[0].Value
	response, err = testClient.Post(testClient.Url("login/refresh"), "", nil)
	if err != nil {
		t.Fatal("Error getting response from server: ", err.Error())
	}
	assert.EqualValues(t, 204, response.StatusCode)

	assert.NotEqual(t, originalToken, response.Cookies()[0].Value, "server should return a different session token")
	assert.NotEqual(t, originalToken, testClient.Jar.Cookies(testClient.url)[0].Value, "client cookiejar should store a different session token")

}
