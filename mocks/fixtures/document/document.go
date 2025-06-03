package documentFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	DocumentFetchResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.DocumentFetchAllRow{
			{
				ID:                    1,
				UserID:                uuid.New(),
				BirthCertificate:      "https://google.com",
				FamilyCard:            "https://google.com",
				UserIdentity:          "https://google.com",
				BeltCertificate:       "https://google.com",
				ElementaryCertificate: "https://google.com",
				MiddleCertificate:     "https://google.com",
				HighCertificate:       "https://google.com",
				BachelorCertificate:   "https://google.com",
				MasterCertificate:     "https://google.com",
			},
		},
	}
)
