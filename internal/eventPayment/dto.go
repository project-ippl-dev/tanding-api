package eventPayment

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type request struct {
	EventID       uuid.UUID   `param:"event"`
	UniqueNumber  int16       `json:"unique_number"`
	Link          string      `json:"link"`
	Total         int32       `json:"total"`
	Registrations []uuid.UUID `json:"registrations"`
}

func (r request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Link, validation.Required, is.URL),
		validation.Field(&r.UniqueNumber, validation.Required),
		validation.Field(&r.Total, validation.Required),
		validation.Field(&r.Registrations, validation.Required),
	)
}

type fetchAllParams struct {
	EventID string                `query:"event_id"`
	Status  db.EventReceiptStatus `query:"status"`
	Clubs   []string              `query:"clubs"`
	Start   string                `query:"start"`
	End     string                `query:"end"`
}

type updateParams struct {
	EventID   uuid.UUID             `param:"event"`
	PaymentID uuid.UUID             `param:"payment"`
	Status    db.EventReceiptStatus `json:"status"`
}

func (u updateParams) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.EventID, validation.Required, is.UUID),
		validation.Field(&u.PaymentID, validation.Required, is.UUID),
		validation.Field(&u.Status, validation.Required, validation.In(
			db.EventReceiptStatusApproved, db.EventReceiptStatusRejected, db.EventReceiptStatusRefund,
		)),
	)
}

type fetchByUserIDParams struct {
	EventID uuid.UUID             `param:"event"`
	Status  db.EventReceiptStatus `query:"status"`
	Start   string                `query:"start"`
	End     string                `query:"end"`
}

type fetchByEventPrivilegeParams struct {
	EventID string                `param:"event"`
	Status  db.EventReceiptStatus `query:"status"`
	Clubs   []string              `query:"clubs"`
	Start   string                `query:"start"`
	End     string                `query:"end"`
}

type response struct {
	fetchAllRow
	ClassEvents []db.ClassEventFetchByPaymentIDRow `json:"class_events"`
}

type cartDetailResponse struct {
	Results      []EventRegistrationFetchCartDetailsRow `json:"results"`
	Event        db.EventFetchForCartRow                `json:"event"`
	UniqueNumber uniqueNumberData                       `json:"unique_number"`
}

type uniqueNumberData struct {
	Number string `json:"number"`
	Time   string `json:"time"`
}

type detailPaymentResponse struct {
	Detail      db.EventPaymentFetchOneForAdminRow     `json:"detail"`
	ClassEvents []EventRegistrationFetchCartDetailsRow `json:"class_events"`
}

type summaryParams struct {
	EventID string                     `param:"event" query:"event_id"`
	Status  db.EventRegistrationStatus `query:"status"`
}

type summaryResponse struct {
	TotalApproved int64 `json:"approved"`
	TotalWaiting  int64 `json:"waiting"`
	TotalRefund   int64 `json:"refund"`
}
