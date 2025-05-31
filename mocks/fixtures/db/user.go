package dbFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

const (
	Name     = "valid"
	Username = "valid"
)

var (
	UserFetchByKeywordRows = []db.UserFetchByKeywordRow{
		{
			ID:       uuid.New(),
			Name:     Name,
			Username: Username,
		},
	}
)
