package club

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type request struct {
	Name      string  `json:"name"`
	Logo      string  `json:"logo"`
	Phone     string  `json:"phone"`
	ShortName string  `json:"short_name"`
	Sports    []sport `json:"sports"`
}

type sport struct {
	SportID string `json:"sport_id"`
}

func (r request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Logo, validation.Required, is.URL),
		validation.Field(&r.Phone, validation.Required),
		validation.Field(&r.ShortName, validation.Required),
		validation.Field(&r.Sports),
	)
}

func (s sport) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.SportID, is.UUID),
	)
}
