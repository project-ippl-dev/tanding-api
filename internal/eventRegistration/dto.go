package eventRegistration

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type RegistrationRequest struct {
	EventID      string       `param:"event"`
	ClassEventID string       `json:"class_event_id"`
	ClubID       string       `json:"club_id"`
	Members      []MemberData `json:"members"`
}

func (r RegistrationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.EventID, validation.Required, is.UUID),
		validation.Field(&r.ClubID, validation.Required, is.UUID),
		validation.Field(&r.ClassEventID, validation.Required, is.UUID),
		validation.Field(&r.Members, validation.Required),
	)
}

type MemberData struct {
	UserID uuid.UUID `json:"user_id"`
}

func (m MemberData) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required, is.UUID),
	)
}

type FetchAllParams struct {
	EventID      string                     `param:"event"`
	ClubID       string                     `query:"club_id"`
	ClassEventID string                     `query:"class_event_id"`
	Status       db.EventRegistrationStatus `query:"status"`
	UserID       string                     `query:"user_id"`
}

type FetchResponse struct {
	FetchAllRow
	Participants []db.EventParticipantFetchByRegistrationIDRow `json:"participants"`
}

type UpdateRegistrationRequest struct {
	EventID      string    `param:"event"`
	RegisterID   uuid.UUID `param:"register"`
	ClassEventID string    `json:"class_event_id"`
}

func (u UpdateRegistrationRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ClassEventID, validation.Required, is.UUID),
	)
}

type SetStatusRequest struct {
	EventID    string    `param:"event"`
	RegisterID uuid.UUID `param:"register"`
}

type FetchParticipantRow struct {
	db.EventRegistrationFetchClubByEventIDRow
	TotalPoint int64                                           `json:"total_point"`
	TotalUser  int64                                           `json:"total_user"`
	Members    []db.EventParticipantFetchByEventIDAndClubIDRow `json:"members"`
}
