package goth

import (
	"crypto/rsa"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"goth/utils"
	"log"
	"os"
	"strconv"
	"time"
)

var SESSION_DATA_KEY string = "sessionData"
var SESSION_PERMISSIONS_KEY string = "sessionPermissions"
var SESSION_USERNAME_KEY string = "sessionUsername"

type JWTProvider struct {
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	alg             string
	validitySeconds int
}

func getPrivateKey() (privateKey *rsa.PrivateKey, err error) {
	var privateKeyData []byte
	privateKeyData, err = os.ReadFile(utils.GetEnvOrDefault("PRIVATE_KEY_PATH", "./id_rsa"))
	if err != nil {
		return
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return
	}
	return
}

func getPublicKey() (publicKey *rsa.PublicKey, err error) {
	var publicKeyData []byte
	publicKeyData, err = os.ReadFile(utils.GetEnvOrDefault("PUBLIC_KEY_PATH", "./id_rsa.pub"))
	if err != nil {
		return
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return
	}
	return
}

func NewJWTProvider(isIssuer bool) (JWTProvider, error) {
	provider := JWTProvider{}
	if isIssuer {
		privateKey, err := getPrivateKey()
		if err != nil {
			return JWTProvider{}, err
		}
		provider.privateKey = privateKey
		provider.publicKey = &privateKey.PublicKey
	} else {
		publicKey, err := getPublicKey()
		if err != nil {
			return provider, err
		}
		provider.publicKey = publicKey
	}

	provider.alg = utils.GetEnvOrDefault("JWT_ALGORITHM", "RS256")
	validitySeconds, err := strconv.Atoi(utils.GetEnvOrDefault("JWT_VALIDITY_SECONDS", "300"))
	if err != nil {
		return provider, err
	}
	provider.validitySeconds = validitySeconds
	return provider, nil
}

type JWTCustomData struct {
	Username string                 `json:"username"`
	Data     map[string]interface{} `json:"data"`
	Roles    []string               `json:"roles"`
}

func (data *JWTCustomData) getPermissions() utils.Set {
	permissionSet := utils.Set{}
	for _, role := range data.Roles {
		if permissions, ok := RolePermissionMap.Get(role); ok {
			permissionSet.Add(permissions...)
		} else {
			log.Printf("[WARNING] missing role mapping for %s", role)
		}
	}
	return permissionSet
}

type jwtStructure struct {
	jwt.RegisteredClaims
	JWTCustomData
}

func (svc JWTProvider) generateToken(customClaims JWTCustomData) (string, error) {
	jwtUUID, _ := uuid.NewV6()
	claims := &jwtStructure{
		jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{time.Now()},
			Issuer:    utils.GetEnv("JWT_ISSUER"),
			ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Second * time.Duration(svc.validitySeconds))},
			ID:        jwtUUID.String(),
		},
		customClaims,
	}

	return jwt.NewWithClaims(jwt.GetSigningMethod(svc.alg), claims).SignedString(svc.privateKey)
}

func (svc JWTProvider) ReadToken(token string) (JWTCustomData, error) {
	//TODO: Does the data return with error? (expiry)
	parsed, err := jwt.ParseWithClaims(
		token,
		&jwtStructure{},
		func(token *jwt.Token) (interface{}, error) { return svc.publicKey, nil },
		jwt.WithValidMethods([]string{svc.alg}),
	)
	claims := parsed.Claims.(*jwtStructure)
	return claims.JWTCustomData, err
}

//TODO: Merge?

func GetPermissions(c *gin.Context) (utils.Set, error) {
	if permissions, ok := c.Get(SESSION_PERMISSIONS_KEY); ok {
		return permissions.(utils.Set), nil
	}
	return nil, fmt.Errorf("missing or invalid permissions in context")
}

func GetSessionData(c *gin.Context) (map[string]interface{}, error) {
	if data, ok := c.Get(SESSION_DATA_KEY); ok {
		return data.(map[string]interface{}), nil
	}
	return nil, fmt.Errorf("missing or invalid data in context")
}

func GetUsername(c *gin.Context) (string, error) {
	if data, ok := c.Get(SESSION_USERNAME_KEY); ok {
		return data.(string), nil
	}
	return "", fmt.Errorf("missing or invalid username in context")
}
