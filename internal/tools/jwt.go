package tools

//go:generate mockgen -source=./jwt.go -destination=../../mocks/tools/jwt_mock.go

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/project-ippl-dev/tanding-api/config"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

// JWT Struct
type JWT struct {
	ID        string
	AccountID int64
	RoleName  db.Role
}

// JWTResponse Struct
type JWTResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"type"`
	ExpiredAt   int64  `json:"expired_at"`
}

type JWTClient interface {
	Auth() echo.MiddlewareFunc
	CreateToken(req JWT) (JWTResponse, error)
	TokenParse(accessToken string) (JWT, error)
	Decode(c echo.Context) JWT
}

type jwtClient struct {
	conf config.JWTConfig
}

func NewJWTClient(jwtConf config.JWTConfig) JWTClient {
	return &jwtClient{
		conf: jwtConf,
	}
}

func (ths *jwtClient) Auth() echo.MiddlewareFunc {
	jwtMiddleware := middleware.JWT([]byte(ths.conf.SecretKey))
	return jwtMiddleware
}

func (ths *jwtClient) CreateToken(request JWT) (JWTResponse, error) {
	expiredAt := time.Now().Add(time.Hour * 24).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = request.ID
	atClaims["account_id"] = request.AccountID
	atClaims["iat"] = time.Now().Unix()
	atClaims["role_name"] = request.RoleName
	atClaims["exp"] = expiredAt
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(ths.conf.SecretKey))
	response := JWTResponse{
		AccessToken: token,
		Type:        "bearer",
		ExpiredAt:   expiredAt,
	}

	return response, err
}

func (ths *jwtClient) TokenParse(accessToken string) (JWT, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(ths.conf.SecretKey), nil
	})
	if err != nil {
		return JWT{}, err
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	accountID := int64(claims["account_id"].(float64))
	roleName := db.Role(claims["role_name"].(string))
	response := JWT{
		ID:        userID,
		AccountID: accountID,
		RoleName:  roleName,
	}

	return response, nil
}

func (ths *jwtClient) Decode(c echo.Context) JWT {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	result := JWT{
		ID:        claims["id"].(string),
		RoleName:  db.Role(claims["role_name"].(string)),
		AccountID: int64(claims["account_id"].(float64)),
	}
	return result
}

func (ths *jwtClient) Middleware() echo.MiddlewareFunc {
	conf := middleware.JWTConfig{
		SigningKey:  []byte(ths.conf.SecretKey),
		TokenLookup: "header:Authorization",
	}
	return middleware.JWTWithConfig(conf)
}
