package class

import (
	"github.com/dytlan/tanding-api/internal/db"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type request struct {
	SportID     string       `json:"sport_id"`
	Name        string       `json:"name"`
	ClassRuleID int64        `json:"class_rule_id"`
	Type        db.ClassType `json:"class_type"`
	MatchType   db.MatchType `json:"match_type"`
}

func (r request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SportID, validation.Required, is.UUID),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.ClassRuleID, validation.Required),
		validation.Field(&r.Type, validation.Required),
	)
}
