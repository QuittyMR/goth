package standalone

import (
	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
	"strings"
)

type postgresService struct {
	connection *pgx.Conn
}

type PostgresUserService postgresService
type PostgresAuthorizationService postgresService

func (svc PostgresAuthorizationService) GetRolePermissionMapping(roleIdentifiers ...string) (mapping map[string][]string) {
	mapping = map[string][]string{}
	query := `select 
				roles.name,
				array_agg(permissions.name) permissions
			from roles 
			inner join roles_permissions rp on rp.role_id = roles.id 
			inner join permissions on permissions.id = rp.permission_id `

	var queryParams []interface{}
	if roleIdentifiers != nil {
		query += "where roles.name in ($1)"
		queryParams = append(queryParams, strings.Join(roleIdentifiers, ","))
	}
	query += "group by roles.name"

	results, err := svc.connection.Query(
		context.Background(),
		query,
		queryParams...,
	)
	if err != nil {
		return nil
	}
	var permissions []string
	var role string
	for results.Next() {
		err := results.Scan(&role, &permissions)
		if err != nil {
			return nil
		}
		mapping[role] = permissions
	}
	return mapping
}

func (svc PostgresAuthorizationService) GetAllPermissions() (permissions []string) {
	query, err := svc.connection.Query(
		context.Background(),
		"select name from permissions",
	)
	if err != nil {
		return nil
	}
	var permission string
	for query.Next() {
		err = query.Scan(&permission)
		if err != nil {
			return nil
		}
		permissions = append(permissions, permission)
	}
	return permissions
}

func NewPostgresUserService(connection *pgx.Conn) PostgresUserService {
	return PostgresUserService{connection: connection}
}

func NewPostgresAuthorizationService(connection *pgx.Conn) PostgresAuthorizationService {
	return PostgresAuthorizationService{connection: connection}
}

func (u PostgresUserService) Get(identifier string) (username, password string, publicData map[string]interface{}, err error) {
	var firstName, lastName string
	var isAdmin bool

	err = u.connection.QueryRow(
		context.Background(),
		"select email, password, first_name, last_name, is_admin from users where email = $1 and is_active",
		identifier,
	).Scan(&username, &password, &firstName, &lastName, &isAdmin)
	if err != nil {
		return
	}
	publicData = map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
	}
	if isAdmin {
		publicData["is_admin"] = isAdmin
	}
	return
}

func (svc PostgresAuthorizationService) GetRolesForUsers(userIdentifiers ...string) (roles []string) {
	if userIdentifiers == nil {
		return
	}
	query, err := svc.connection.Query(
		context.Background(),
		`select 
				roles.name, 
				roles.type 
			from users_roles ur 
			left join roles on ur.role_id = roles.id 
			left join users on users.id = ur.user_id 
			where email in ($1)`,
		strings.Join(userIdentifiers, ","),
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

func (svc PostgresAuthorizationService) GetAllRoles() (roles []string) {
	query, err := svc.connection.Query(
		context.Background(),
		"select name from roles",
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
