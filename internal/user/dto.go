package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type searchParams struct {
	Keyword string `query:"keyword"`
	Limit   int32  `query:"limit"`
}

type updateBasicInformationParams struct {
	UserID         string `param:"uuid"`
	Name           string `json:"name"`
	BornAt         string `json:"born_at"`
	BornOn         string `json:"born_on"`
	IdentityNumber string `json:"identity_number"`
	Phone          string `json:"phone"`
	Photo          string `json:"photo"`
	Gender         string `json:"gender"`
	About          string `json:"about"`
}

func (u updateBasicInformationParams) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Photo, is.URL),
		validation.Field(&u.Phone, validation.Length(10, 15)),
		validation.Field(&u.IdentityNumber, is.Int, validation.Length(16, 16)),
		validation.Field(&u.Gender, validation.Required, validation.In("male", "female")),
	)
}

type basicInformationResponse struct {
	db.UserFetchBasicInformationRow
	EventPrivilege bool
}
