package dbFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

var (
	EventPrivilegeFetchAllRows = []db.EventPrivilegeFetchAllRow{
		{
			ID:     1,
			Role:   db.EventRoleContributor,
			UserID: uuid.New(),
			Name:   "name",
		},
	}
)
