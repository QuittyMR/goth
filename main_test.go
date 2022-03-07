package gauth

import (
	"fmt"
	"gauth/internal"
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
	userService := internal.NewUserService(connection)
	testServer = httptest.NewServer(gauthServer)
	//gauthServer := GetAuthGroup(testServer, userService, )
	defer testServer.Close()
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	client := NewClient()

	body := strings.NewReader(`{"email":"t1@example.com", "password":"1234"}`)
	_, _ = client.Post(
		fmt.Sprintf("%s/%s", testServer.URL, "login/basic"),
		"application/json",
		body,
	)
	_, _ = client.Get(fmt.Sprintf("%s/%s", testServer.URL, "alive"))
}
