package certificateFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/certificate"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"time"
)

var (
	CertificateResponse = certificate.Response{
		Certificate: db.CertificateFetchOneRow{
			ID:        uuid.New(),
			Name:      "name",
			EventID:   uuid.New(),
			UserID:    uuid.New(),
			EventName: "event-name",
			RewardAs:  "1st",
			CreatedAt: time.Now(),
		},
		Recipient: "name@email.com",
	}

	CertificateFetchResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.CertificateFetchAllByUserIDRow{
			{
				ID:        uuid.New(),
				Name:      "name",
				EventName: "event-name",
				Thumbnail: "https://google.com",
				RewardAs:  "1st",
				CreatedAt: time.Now(),
			},
		},
	}
)
