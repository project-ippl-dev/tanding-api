package eventRegistrationFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/eventRegistration"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	EventRegistrationID = uuid.New()
)

var (
	EventRegistrationFetchAllResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []eventRegistration.FetchResponse{
			{
				FetchAllRow: eventRegistration.FetchAllRow{
					ID:           uuid.New(),
					ClassEventID: uuid.New(),
					ClassName:    "class name",
					ClubID:       uuid.New(),
					ClubName:     "club name",
					Status:       db.EventRegistrationStatusApproved,
					Price:        1000,
				},
				Participants: []db.EventParticipantFetchByRegistrationIDRow{
					{
						ID:                  1,
						EventRegistrationID: uuid.New(),
						UserID:              uuid.New(),
						Name:                "name",
						ClassName:           "class name",
					},
				},
			},
		},
	}

	EventRegistrationFetchParticipantRows = []eventRegistration.FetchParticipantRow{
		{
			EventRegistrationFetchClubByEventIDRow: db.EventRegistrationFetchClubByEventIDRow{
				ID:   uuid.New(),
				Name: "name",
			},
			TotalPoint: 1,
			TotalUser:  1,
			Members: []db.EventParticipantFetchByEventIDAndClubIDRow{
				{
					ID:                  1,
					EventRegistrationID: uuid.New(),
					UserID:              uuid.New(),
					Name:                "name",
					ClassName:           "class name",
				},
			},
		},
	}
)
