package standalone

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/publicsuffix"
	"goth"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

var testServer *httptest.Server
var testRecorder *httptest.ResponseRecorder
var testContext *gin.Context

type GothTestClient struct {
	*http.Client
	url *url.URL
	t   *testing.T
}

func NewClient(t *testing.T) GothTestClient {
	uri, _ := url.ParseRequestURI(testServer.URL)
	client := GothTestClient{testServer.Client(), uri, t}
	client.Jar = newCookieJar()
	return client
}

func newCookieJar() *cookiejar.Jar {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return cookieJar
}

func (client GothTestClient) Url(path string) string {
	return fmt.Sprintf("%s/auth/%s", client.url, path)
}

func (client GothTestClient) GetToken() goth.JWTCustomData {
	cookies := client.Jar.Cookies(client.url)
	if cookies == nil {
		client.t.Fatal("cannot find cookies")
	}
	token, err := jwtService.ReadToken(cookies[0].Value)
	if err != nil {
		client.t.Fatal("error getting data from token:", err.Error())
	}
	return token
}

func GetJSONBody(response *http.Response) (mapping map[string]interface{}, err error) {
	contentLength, err := strconv.Atoi(response.Header["Content-Length"][0])
	if err != nil {
		return
	}
	resData := make([]byte, contentLength, contentLength)
	_, err = response.Body.Read(resData)
	if err != nil {
		return
	}
	err = json.Unmarshal(resData, &mapping)
	if err != nil {
		return
	}
	return
}

type mockUserService struct{}

func (svc mockUserService) Get(identifier string) (password string, userData map[string]interface{}) {
	return "1234", map[string]interface{}{
		"username":      identifier,
		"first_name":    "Test",
		"last_name":     "McTesty",
		"random_data":   rand.Int(),
		"compound_data": map[string]string{"data": "nested"},
	}
}

type mockAuthorizationService struct {
	rolePermissionMapping map[string][]string
}

func newMockAuthorizationService() mockAuthorizationService {
	svc := mockAuthorizationService{rolePermissionMapping: map[string][]string{
		"manager": {"get_roles", "get_permissions", "overlapping_permission"},
		"member":  {"get_roles_self", "overlapping_permission"},
	}}
	return svc
}

func (svc mockAuthorizationService) GetRolesForUsers(userIdentifiers ...string) []string {
	panic("implement me")
}

func (svc mockAuthorizationService) GetAllRoles() []string {
	panic("implement me")
}

func (svc mockAuthorizationService) GetRolePermissionMapping(roleIdentifiers ...string) map[string][]string {
	mapping := map[string][]string{}
	for _, role := range roleIdentifiers {
		mapping[role] = svc.rolePermissionMapping[role]
	}
	return mapping
}

func (svc mockAuthorizationService) GetAllPermissions() []string {
	panic("implement me")
}
