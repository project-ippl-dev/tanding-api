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

type participantReq struct {
	ClubID       string            `param:"club"`
	Participants []participantData `json:"participants"`
}

func (p participantReq) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Participants, validation.Required),
		validation.Field(&p.ClubID, validation.Required, is.UUID),
	)
}

type participantData struct {
	UserID  string `json:"user_id"`
	SportID string `json:"sport_id"`
}

func (p participantData) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserID, validation.Required, is.UUID),
		validation.Field(&p.SportID, validation.Required, is.UUID),
	)
}

type updateJoinApprovalArgs struct {
	ClubID     string `param:"club"`
	Status     *bool  `json:"status"`
	ApprovalID int64  `param:"approval"`
}

func (u updateJoinApprovalArgs) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ClubID, validation.Required, is.UUID),
		validation.Field(&u.ApprovalID, validation.Required),
	)
}

type updateInviteApprovalArgs struct {
	Status     *bool `json:"status"`
	ApprovalID int64 `param:"approval"`
}

func (u updateInviteApprovalArgs) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ApprovalID, validation.Required),
	)
}

type joinParam struct {
	ClubID  string `param:"club"`
	SportID string `json:"sport_id"`
}

func (j joinParam) Validate() error {
	return validation.ValidateStruct(&j,
		validation.Field(&j.SportID, is.UUID),
		validation.Field(&j.ClubID, validation.Required, is.UUID),
	)
}
