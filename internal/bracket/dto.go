package bracket

import (
	"context"
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type generateParams struct {
	ClassEventID uuid.UUID `param:"class"`
}

type orderRoundDownResponse struct {
	OrderBracketFetchByClassEventIDRow
	Iteration int16 `json:"iteration"`
}

type updateLockParams struct {
	ClassEventID uuid.UUID               `param:"class"`
	Status       *bool                   `json:"status"`
	Participants []participantLockParams `json:"participants"`
}

func (u updateLockParams) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Participants, validation.When(*u.Status, validation.Required)),
	)
}

type participantLockParams struct {
	EventRegistrationID string `json:"event_registration_id"`
	Iteration           int16  `json:"iteration"`
}

func (p participantLockParams) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.EventRegistrationID, is.UUID),
		validation.Field(&p.Iteration, validation.When(p.EventRegistrationID != "", validation.Required)),
	)
}

type updateGenerateParams struct {
	ClassEventID uuid.UUID `param:"class"`
}

type fetchOneResponse struct {
	Message        string                       `json:"message"`
	Data           interface{}                  `json:"data"`
	GenerateStatus *bool                        `json:"generate_status"`
	LockStatus     *bool                        `json:"lock_status"`
	MatchType      db.MatchType                 `json:"match_type"`
	LockScore      *bool                        `json:"lock_score"`
	Summary        []RankFetchByClassEventIDRow `json:"summary"`
}

type roundDownResponse struct {
	Message   string       `json:"message"`
	Data      interface{}  `json:"data"`
	MatchType db.MatchType `json:"match_type"`
}

type storeBracketParams struct {
	Ctx               context.Context
	Tx                *sql.Tx
	TxQuery           *db.Queries
	ClassEvent        db.ClassEventFetchOneRow
	RegistrationCount int
	MatchIndex        int
}

type matchIndexData struct {
	Title string     `json:"title"`
	Seeds []seedData `json:"seeds"`
}

type seedData struct {
	ID         uuid.UUID                    `json:"id"`
	EventTurn  int16                        `json:"event_turn"`
	MatchOrder int16                        `json:"match_order"`
	IsActive   int16                        `json:"is_active"`
	IsScore    bool                         `json:"is_score"`
	Teams      []bracketParticipantResponse `json:"teams"`
}

type bracketParticipantResponse struct {
	bracketParticipantFetchByEventBracketIDRow
	ClubLogo string `json:"club_logo"`
	Score    bracketParticipantScoreParams
}

type bracketParticipantScoreParams struct {
	Round1 int16 `json:"round1"`
	Round2 int16 `json:"round2"`
	Round3 int16 `json:"round3"`
	Extra  int16 `json:"extra"`
	Total  int16 `json:"total"`
}

func (s seedData) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required),
		validation.Field(&s.MatchOrder, validation.Required),
		validation.Field(&s.Teams, validation.Required),
	)
}

func (b bracketParticipantFetchByEventBracketIDRow) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.EventRegistrationID, validation.Required),
	)
}

type updateSingleLockParams struct {
	EventID      uuid.UUID      `param:"event"`
	ClassEventID uuid.UUID      `param:"class"`
	Data         matchIndexData `json:"data"`
	Status       *bool          `json:"status"`
}

func (u updateSingleLockParams) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Data, validation.Required),
	)
}

type fetchOneOrderResponse struct {
	OrderBracketFetchByClassEventIDRow
	Scores db.OrderScoreFetchOneByBracketIDRow `json:"scores"`
}
