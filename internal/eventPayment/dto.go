package eventPayment

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Request struct {
	EventID       uuid.UUID   `param:"event"`
	UniqueNumber  int16       `json:"unique_number"`
	Link          string      `json:"link"`
	Total         int32       `json:"total"`
	Registrations []uuid.UUID `json:"registrations"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Link, validation.Required, is.URL),
		validation.Field(&r.UniqueNumber, validation.Required),
		validation.Field(&r.Total, validation.Required),
		validation.Field(&r.Registrations, validation.Required),
	)
}

type FetchAllParams struct {
	EventID string                `query:"event_id"`
	Status  db.EventReceiptStatus `query:"status"`
	Clubs   []string              `query:"clubs"`
	Start   string                `query:"start"`
	End     string                `query:"end"`
}

type UpdateParams struct {
	EventID   uuid.UUID             `param:"event"`
	PaymentID uuid.UUID             `param:"payment"`
	Status    db.EventReceiptStatus `json:"status"`
}

func (u UpdateParams) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.EventID, validation.Required, is.UUID),
		validation.Field(&u.PaymentID, validation.Required, is.UUID),
		validation.Field(&u.Status, validation.Required, validation.In(
			db.EventReceiptStatusApproved, db.EventReceiptStatusRejected, db.EventReceiptStatusRefund,
		)),
	)
}

type FetchByUserIDParams struct {
	EventID uuid.UUID             `param:"event"`
	Status  db.EventReceiptStatus `query:"status"`
	Start   string                `query:"start"`
	End     string                `query:"end"`
}

type FetchByEventPrivilegeParams struct {
	EventID string                `param:"event"`
	Status  db.EventReceiptStatus `query:"status"`
	Clubs   []string              `query:"clubs"`
	Start   string                `query:"start"`
	End     string                `query:"end"`
}

type Response struct {
	FetchAllRow
	ClassEvents []db.ClassEventFetchByPaymentIDRow `json:"class_events"`
}

type CartDetailResponse struct {
	Results      []EventRegistrationFetchCartDetailsRow `json:"results"`
	Event        db.EventFetchForCartRow                `json:"event"`
	UniqueNumber UniqueNumberData                       `json:"unique_number"`
}

type UniqueNumberData struct {
	Number string `json:"number"`
	Time   string `json:"time"`
}

type DetailPaymentResponse struct {
	Detail      db.EventPaymentFetchOneForAdminRow     `json:"detail"`
	ClassEvents []EventRegistrationFetchCartDetailsRow `json:"class_events"`
}

type SummaryParams struct {
	EventID string                     `param:"event" query:"event_id"`
	Status  db.EventRegistrationStatus `query:"status"`
}

type SummaryResponse struct {
	TotalApproved int64 `json:"approved"`
	TotalWaiting  int64 `json:"waiting"`
	TotalRefund   int64 `json:"refund"`
}
