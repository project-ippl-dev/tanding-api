package club

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Request struct {
	Name      string  `json:"name"`
	Logo      string  `json:"logo"`
	Phone     string  `json:"phone"`
	ShortName string  `json:"short_name"`
	Sports    []sport `json:"sports"`
}

type sport struct {
	SportID string `json:"sport_id"`
}

func (r Request) Validate() error {
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

type ParticipantReq struct {
	ClubID       string            `param:"club"`
	Participants []participantData `json:"participants"`
}

func (p ParticipantReq) Validate() error {
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

type UpdateJoinApprovalArgs struct {
	ClubID     string `param:"club"`
	Status     *bool  `json:"status"`
	ApprovalID int64  `param:"approval"`
}

func (u UpdateJoinApprovalArgs) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ClubID, validation.Required, is.UUID),
		validation.Field(&u.ApprovalID, validation.Required),
	)
}

type UpdateInviteApprovalArgs struct {
	Status     *bool `json:"status"`
	ApprovalID int64 `param:"approval"`
}

func (u UpdateInviteApprovalArgs) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ApprovalID, validation.Required),
	)
}

type JoinParam struct {
	ClubID  string `param:"club"`
	SportID string `json:"sport_id"`
}

func (j JoinParam) Validate() error {
	return validation.ValidateStruct(&j,
		validation.Field(&j.SportID, is.UUID),
		validation.Field(&j.ClubID, validation.Required, is.UUID),
	)
}
