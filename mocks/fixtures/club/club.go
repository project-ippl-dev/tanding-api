package clubFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/club"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	ClubFetchParticipantResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: club.FetchParticipantResponse{
			TotalPoint: 10,
			Participants: []club.FetchParticipantRow{
				{
					ParticipantFetchAllRow: club.ParticipantFetchAllRow{
						ID:             1,
						UserID:         uuid.New(),
						Name:           "name",
						CanParticipate: true,
					},
					Point: 10,
				},
			},
		},
	}

	ClubFetchAllResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []club.FetchAllRow{
			{
				ID:    uuid.New(),
				Name:  "name",
				Logo:  "https://google.com",
				Owner: "owner",
			},
		},
	}

	ClubFetchOneResponse = club.FetchOneResponse{
		ClubFetchOneRow: db.ClubFetchOneRow{
			ID:        uuid.New(),
			Name:      "name",
			Logo:      "https://google.com",
			ShortName: "shr",
			Owner:     "owner",
		},
		Privilege: true,
		Joined:    true,
	}
)
