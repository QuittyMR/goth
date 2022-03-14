package gauth

import (
	"gauth/standalone"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	//TODO: Mock DB
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/test_auth")
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	defer connection.Close(context.Background())
	jwtService, err := NewJWTProvider()
	if err != nil {
		panic(err)
	}
	testRecorder = httptest.NewRecorder()
	_, ginServer := gin.CreateTestContext(testRecorder)

	_ = GetAuthGroup(
		ginServer,
		standalone.NewPostgresUserService(connection),
		standalone.NewPostgresAuthorizationService(connection),
		jwtService,
	)
	testServer = httptest.NewTLSServer(ginServer)
	defer testServer.Close()

	os.Exit(m.Run())
}

func TestAuthService_BasicLogin(t *testing.T) {
	client := testServer.Client()

	request := httptest.NewRequest(
		"POST",
		"authenticate/basic",
		strings.NewReader(`{"email":"t1@example.com", "password":"1234"}`),
	)
	response, err := client.Do(request)
	if err != nil {
		return
	}
	recorderResponse := testRecorder.Result()
	assert.Equal(t, response, recorderResponse)
}
