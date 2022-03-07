package internal

import (
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

//TODO: once? singleton?
type postgresService struct {
	connection *pgx.Conn
}

type PostgresUserService postgresService
type PostgresRoleService postgresService
type PostgresPermissionService postgresService

func (svc PostgresPermissionService) GetForRole(roleIdentifier string) []string {
	panic("implement me")
}

func (svc PostgresPermissionService) GetAll() []string {
	panic("implement me")
}

func NewUserService(connection *pgx.Conn) PostgresUserService {
	return PostgresUserService{connection: connection}
}

func NewPostgresRoleService(connection *pgx.Conn) PostgresRoleService {
	return PostgresRoleService{connection: connection}
}

func NewPermissionservice(connection *pgx.Conn) PostgresPermissionService {
	return PostgresPermissionService{connection: connection}
}

func (u PostgresUserService) Get(identifier string) (password string, publicData map[string]interface{}) {
	var email, name string

	err := u.connection.QueryRow(
		context.Background(),
		"select email, password, name from users where email = $1 and is_active",
		identifier,
	).Scan(&email, &password, &name)
	if err != nil {
		return "", nil
	}
	return password, map[string]interface{}{"email": email, "name": name}
}

func (svc PostgresRoleService) GetForUser(userIdentifier string) (roles []string) {
	query, err := svc.connection.Query(
		context.Background(),
		`select 
				roles.name, 
				roles.type 
			from users_roles ur 
			left join roles on ur.role_id = roles.id 
			left join users on users.id = ur.user_id 
			where email = $1`,
		userIdentifier,
	)
	if err != nil {
		return nil
	}
	var role string
	var roleType int
	for query.Next() {
		err := query.Scan(&role, &roleType)
		if err != nil {
			return nil
		}
		roles = append(roles, role)
	}
	return roles
}

func (svc PostgresRoleService) GetAll() (roles []string) {
	query, err := svc.connection.Query(
		context.Background(),
		"select name, type from roles",
	)
	if err != nil {
		return nil
	}
	var role string
	for query.Next() {
		err = query.Scan(&role)
		if err != nil {
			return nil
		}
		roles = append(roles, role)
	}
	return roles
}
