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

	USerFetchOneRow = db.UserFetchOneRow{
		ID:        uuid.New(),
		AccountID: 1,
		Name:      "name",
		Photo:     "https://google.com",
		Username:  "name@email.com",
		Type:      db.AccountTypeGoogle,
		Role:      db.RoleUser,
	}
)
