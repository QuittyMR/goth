package gauth

import (
	"context"
	"fmt"
	"gauth/login"
	"github.com/jackc/pgx/v4"
	"log"
	"testing"
)

type UserService struct {
	connection *pgx.Conn
}

func (u UserService) GetPassword(email string) (password string) {
	err := u.connection.QueryRow(
		context.Background(),
		"select password from users where email = $1",
		email,
	).Scan(&password)
	if err != nil {
		fmt.Errorf("error retrieving password for %s: %s", email, err.Error())
	}
	return
}

func (u UserService) SetPassword(email, newPassword string) error {
	res, err := u.connection.Exec(context.Background(),
		"update users set password = $1 where email = $2",
		newPassword,
		email,
	)
	if err != nil {
		return err
	}
	log.Println(res.String())
	return nil
}

func (u UserService) SetActive(email string, isActive bool) error {
	_, err := u.connection.Query(context.Background(),
		"update users set password = ? where email = ?",
		[]interface{}{isActive, email},
	)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(connection *pgx.Conn) UserService {
	return UserService{connection: connection}
}

func TestNew(t *testing.T) {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/test_auth")
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	defer connection.Close(context.Background())
	var userService login.UserService = NewUserService(connection)

	server := New(&userService)
	_ = server.Run()
}
