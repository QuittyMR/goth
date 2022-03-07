package login

import (
	"crypto/rsa"
	"gauth/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
	"strconv"
	"time"
)

func (svc JWTService) JWTMiddlware(context *gin.Context) {
	cookie, _ := context.Cookie("session")
	if cookie == "" {
		context.Next()
		return
	}
	token, err := svc.readToken(cookie)
	if err != nil {
		log.Printf("Invalid token received: %s", err.Error())
	} else {
		context.Set("session", token.Data)
	}
	context.Next()
}

type JWTService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	alg        string
}

func getPrivateKey() (key *rsa.PrivateKey, err error) {
	var privateKeyData []byte
	privateKeyData, err = os.ReadFile(utils.GetEnv("PRIVATE_KEY_PATH"))
	if err != nil {
		return
	}
	key, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return
	}
	return key, nil
}

func NewJWTService() (JWTService, error) {
	key, err := getPrivateKey()
	if err != nil {
		return JWTService{}, err
	}
	return JWTService{
		privateKey: key,
		publicKey:  &key.PublicKey,
		alg:        utils.GetEnvOrDefault("JWT_ALGORITHM", "RS256"),
	}, nil
}

type JWTCustomData struct {
	Data map[string]interface{} `json:"data"`
}

type jwtStructure struct {
	*jwt.RegisteredClaims
	JWTCustomData
}

func (svc JWTService) generateToken(customClaims JWTCustomData) (string, error) {
	jwtTimeout, err := strconv.Atoi(utils.GetEnvOrDefault("JWT_EXPIRATION_MINUTES", "5"))
	if err != nil {
		log.Fatalf("JWT error: expiration set incorrectly: %s", err.Error())
	}
	claims := &jwtStructure{
		&jwt.RegisteredClaims{
			IssuedAt:  &jwt.NumericDate{time.Now()},
			Issuer:    utils.GetEnv("JWT_ISSUER"),
			ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Minute * time.Duration(jwtTimeout))},
		},
		customClaims,
	}

	return jwt.NewWithClaims(jwt.GetSigningMethod(svc.alg), claims).SignedString(svc.privateKey)
}

func (svc JWTService) readToken(token string) (*jwtStructure, error) {
	parsed, err := jwt.ParseWithClaims(
		token,
		&jwtStructure{},
		func(token *jwt.Token) (interface{}, error) { return svc.publicKey, nil },
		jwt.WithValidMethods([]string{svc.alg}),
	)
	if err != nil {
		return nil, err
	}
	claims := parsed.Claims.(*jwtStructure)
	return claims, nil
}
