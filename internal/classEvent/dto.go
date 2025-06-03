package classEvent

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Request struct {
	Data []RequestData `json:"data"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Data, validation.Required),
	)
}

type RequestData struct {
	ClassID string `json:"class_id"`
	Price   int32  `json:"price"`
}

func (r RequestData) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ClassID, is.UUID, validation.Required),
		validation.Field(&r.Price, validation.Min(0)),
	)
}

type updateReq struct {
	Price int32 `json:"price"`
}

func (u updateReq) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Price, validation.Min(0)),
	)
}
