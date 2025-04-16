package tools

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

// JWTAuth representation of Authentication Middleware with JWT
func JWTAuth() echo.MiddlewareFunc {
	jwtMiddleware := middleware.JWT([]byte(config.Configuration().JWT.SecretKey))
	return jwtMiddleware
}

// JWTCreateToken representation generate JWT Token
func JWTCreateToken(request JWT) (JWTResponse, error) {
	JWTSecretKey := config.Configuration().JWT.SecretKey
	expiredAt := time.Now().Add(time.Hour * 24).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = request.ID
	atClaims["account_id"] = request.AccountID
	atClaims["iat"] = time.Now().Unix()
	atClaims["role_name"] = request.RoleName
	atClaims["exp"] = expiredAt
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(JWTSecretKey))
	response := JWTResponse{
		AccessToken: token,
		Type:        "bearer",
		ExpiredAt:   expiredAt,
	}

	return response, err
}

// JWTTokenParse represent decode JWT token string to userID
func JWTTokenParse(accessToken string) (JWT, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Configuration().JWT.SecretKey), nil
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

// JWTDecode represent decode JWT based on headers name Authorization
func JWTDecode(c echo.Context) JWT {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	result := JWT{
		ID:        claims["id"].(string),
		RoleName:  db.Role(claims["role_name"].(string)),
		AccountID: int64(claims["account_id"].(float64)),
	}
	return result
}

func JWTMiddleware() echo.MiddlewareFunc {
	conf := middleware.JWTConfig{
		SigningKey:  []byte(config.Configuration().JWT.SecretKey),
		TokenLookup: "header:Authorization,param:token",
	}
	return middleware.JWTWithConfig(conf)
}
