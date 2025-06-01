package jwtFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	DecodedJWT = tools.JWT{
		ID:        uuid.NewString(),
		AccountID: 0,
		RoleName:  db.RoleUser,
	}
)
