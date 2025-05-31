package class

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Request struct {
	SportID     string       `json:"sport_id"`
	Name        string       `json:"name"`
	ClassRuleID int64        `json:"class_rule_id"`
	Type        db.ClassType `json:"class_type"`
	MatchType   db.MatchType `json:"match_type"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.SportID, validation.Required, is.UUID),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.ClassRuleID, validation.Required),
		validation.Field(&r.Type, validation.Required),
	)
}
