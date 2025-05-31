package userFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"github.com/project-ippl-dev/tanding-api/internal/user"
	"time"
)

var (
	UserFetchBasicInformationResponse = user.BasicInformationResponse{
		UserFetchBasicInformationRow: db.UserFetchBasicInformationRow{
			ID:             uuid.New(),
			Name:           "valid",
			BornAt:         "born_at",
			IdentityNumber: "identity_number",
			Phone:          "1234567890",
			Gender:         "male",
			Status:         true,
			CanParticipate: true,
		},
		EventPrivilege: true,
	}

	UserFetchAllResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.UserFetchAllRow{
			{
				ID:        uuid.New(),
				Name:      "name",
				Role:      "role",
				CreatedAt: time.Now(),
			},
		},
	}

	UserFetchLastLoginResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.LoginDetailFetchAllRow{
			{
				Name:      "name",
				LastLogin: time.Now(),
			},
		},
	}
)
