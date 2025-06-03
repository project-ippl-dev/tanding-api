package eventFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/event"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"time"
)

var (
	EventID = uuid.New()
)

var (
	EventFetchAllResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []event.Response{
			{
				ID:          uuid.New(),
				UserID:      uuid.New(),
				UserName:    "username",
				UserImage:   "https://google.com",
				Type:        db.EventTypeCompetition,
				Name:        "name",
				Description: "description",
				PrizePool:   "prizepool",
				Location:    "location",
				Province:    "online",
				City:        "online",
				Thumbnail:   "https://google.com",
				SportID:     uuid.New(),
				SportName:   "sport-name",
			},
		},
	}

	EventFetchOneResponse = event.ResponseFetchOne{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		UserName:     "username",
		UserImage:    "https://google.com",
		Type:         db.EventTypeCompetition,
		Name:         "name",
		Description:  "description",
		PrizePool:    "prizepool",
		Location:     "location",
		Province:     "province",
		City:         "city",
		Thumbnail:    "https://google.com",
		StartDate:    "2025-06-02",
		EndDate:      "2025-06-02",
		Deadline:     time.Now().Add(24 * time.Hour),
		SportID:      uuid.New(),
		SportName:    "sport-name",
		Rules:        "rules",
		ProposalLink: "https://google.com",
		Status:       true,
		Quota:        1000,
		Open:         "open",
		Remark:       "remark",
	}

	EventFetchInfiniteResponse = event.ResponseInfinite{
		Message: "fetch infinite scroll for event success",
		Data: []event.ResponseInfiniteRow{
			{
				ID:          uuid.New(),
				UserID:      uuid.New(),
				UserName:    "username",
				Type:        db.EventTypeCompetition,
				Name:        "name",
				Description: "description",
				PrizePool:   "prizepool",
				Location:    "location",
				Province:    "province",
				City:        "city",
				Thumbnail:   "https://google.com",
			},
		},
		TotalItem: 1,
	}
)
