package auth

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type registerRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (r registerRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Phone, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100)),
		validation.Field(&r.ConfirmPassword, validation.Required, validation.Length(8, 100), validation.By(tools.StringEquals(r.Password))),
	)
}

type sendMailRequest struct {
	BodyParam bodyParam
	BodyPath  string
}

type bodyParam struct {
	Name string
	Path string
	Host *string
}

type loginReq struct {
	Username string
	Password string
}

func (l loginReq) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Username, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required, validation.Length(8, 100)),
	)
}

type loginResponse struct {
	Data           loginDataResponse `json:"profile"`
	Token          tools.JWTResponse `json:"token"`
	Role           db.Role           `json:"role"`
	Privileges     interface{}       `json:"privileges"`
	Clubs          interface{}       `json:"clubs"`
	CanParticipate *bool             `json:"can_participate"`
	UserID         uuid.UUID         `json:"user_id"`
	Owners         interface{}       `json:"owners"`
}

type loginDataResponse struct {
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

// profileData is binding response for all callback
type profileData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	Name     string `json:"name"`
	Photo    string `json:"picture"`
}

// profileFacebook is binding response from facebook callback
type profileFacebook struct {
	ID       string          `json:"id"`
	Email    string          `json:"email"`
	Verified bool            `json:"verified_email"`
	Name     string          `json:"name"`
	Picture  pictureFacebook `json:"picture"`
}

type pictureFacebook struct {
	Data pictureDataFacebook `json:"data"`
}

type pictureDataFacebook struct {
	URL string `json:"url"`
}

type storeFromCallbackResponse struct {
	Name      string
	UserID    uuid.UUID
	AccountID int64
	Photo     string
}

type resetRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (r resetRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100)),
		validation.Field(&r.ConfirmPassword, validation.Required, validation.Length(8, 100)),
	)
}

type bindingRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Token           string `json:"token"`
}

func (b bindingRequest) Validate(kind string) error {
	switch kind {
	case string(db.AccountTypeManual):
		return validation.ValidateStruct(&b,
			validation.Field(&b.Email, validation.Required, is.Email),
			validation.Field(&b.Password, validation.Required, validation.Length(8, 100)),
			validation.Field(&b.ConfirmPassword, validation.Required, validation.Length(8, 100), validation.By(tools.StringEquals(b.Password))),
		)
	case string(db.AccountTypeFacebook), string(db.AccountTypeGoogle):
		return validation.ValidateStruct(&b,
			validation.Field(&b.Token, validation.Required),
		)
	}
	return fmt.Errorf("error type, only validating request for manual, google, or facebook")
}

type bindingParams struct {
	Kind    string
	Request bindingRequest
	Decoded tools.JWT
	Host    string
}

type userIPDetail struct {
	Country string `json:"country"`
	Region  string `json:"region"`
}
