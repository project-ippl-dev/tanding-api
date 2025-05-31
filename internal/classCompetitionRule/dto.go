package classCompetitionRule

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Request struct {
	Name   string `json:"name"`
	Male   int16  `json:"male"`
	Female int16  `json:"female"`
	Total  int16  `json:"total"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Male, validation.Min(0), validation.Max(20)),
		validation.Field(&r.Female, validation.Min(0), validation.Max(20)),
		validation.Field(&r.Total, validation.Required),
	)
}
