package authFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/auth"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"time"
)

var (
	AuthLoginResponse = auth.LoginResponse{
		Data: auth.LoginDataResponse{
			Name:  "name",
			Photo: "https://google.com",
		},
		Token: tools.JWTResponse{
			AccessToken: "valid-token",
			Type:        "bearer",
			ExpiredAt:   time.Now().Add(time.Hour * 72).Unix(),
		},
		Role:   "user",
		UserID: uuid.New(),
	}
)
