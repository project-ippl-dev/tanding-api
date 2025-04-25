package score

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type orderStoreOrUpdateParams struct {
	OrderBracketID uuid.UUID `param:"bracket"`
	Round1         int16     `json:"round_1"`
	Round2         int16     `json:"round_2"`
	Round3         int16     `json:"round_3"`
	Extra          int16     `json:"extra"`
	Total          int16     `json:"total"`
}

func (o orderStoreOrUpdateParams) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Total, validation.Min(0)),
	)
}

type fetchOneParams struct {
	EventID   uuid.UUID `param:"event"`
	BracketID uuid.UUID `param:"bracket"`
}

type lockParams struct {
	EventID      uuid.UUID `param:"event"`
	ClassEventID uuid.UUID `param:"class"`
	Status       *bool     `json:"status"`
}

type singleStoreOrUpdateParams struct {
	EventBracketID uuid.UUID `param:"bracket"`
	HomeRound1     int16     `json:"home_round1"`
	HomeRound2     int16     `json:"home_round2"`
	HomeRound3     int16     `json:"home_round3"`
	HomeExtra      int16     `json:"home_extra"`
	HomeTotal      int16     `json:"home_total"`
	AwayRound1     int16     `json:"away_round1"`
	AwayRound2     int16     `json:"away_round2"`
	AwayRound3     int16     `json:"away_round3"`
	AwayExtra      int16     `json:"away_extra"`
	AwayTotal      int16     `json:"away_total"`
}

func (o singleStoreOrUpdateParams) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.HomeTotal, validation.Min(0)),
		validation.Field(&o.AwayTotal, validation.Min(0)),
	)
}
