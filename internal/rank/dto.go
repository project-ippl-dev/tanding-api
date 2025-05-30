package rank

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type storeCertificateByRegistrationIDParams struct {
	EventRegistrationID uuid.UUID
	Rank                int16
	EventID             uuid.UUID
	ClassEventID        uuid.UUID
	ClassName           string
}

type storeCertificateExcludeRegistrationIDParams struct {
	EventRegistrationID []uuid.UUID
	Rank                int16
	EventID             uuid.UUID
	ClassEventID        uuid.UUID
	ClassName           string
}

type FetchByClubIDResponse struct {
	TotalPoint   int64                             `json:"total_point"`
	Participants []db.RankFetchAllPointByClubIDRow `json:"participants"`
}

type RankParams struct {
	SportID string `query:"sport_id"`
}

func (r RankParams) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SportID, is.UUID),
	)
}
