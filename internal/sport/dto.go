package sport

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type request struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	SportType   db.SportType `json:"sport_type"`
	Thumbnail   string       `json:"thumbnail"`
}

func (r request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.SportType, validation.Required),
		validation.Field(&r.Thumbnail, is.URL),
	)
}

type fetchAllQueryParams struct {
	Keyword  string       `query:"keyword"`
	Category db.SportType `query:"category"`
}
