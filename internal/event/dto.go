package event

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Request struct {
	Name         string       `json:"name"`
	Type         db.EventType `json:"type"`
	Description  string       `json:"description"`
	PrizePool    string       `json:"prize_pool"`
	Location     string       `json:"location"`
	Province     string       `json:"province"`
	City         string       `json:"city"`
	Thumbnail    string       `json:"thumbnail"`
	StartDate    string       `json:"start_date"`
	EndDate      string       `json:"end_date"`
	Deadline     string       `json:"deadline"`
	SportID      string       `json:"sport_id"`
	Rules        string       `json:"rules"`
	Open         string       `json:"open"`
	Quota        int32        `json:"quota"`
	ProposalLink string       `json:"proposal_link"`
}

func (r Request) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Type, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.PrizePool, validation.Required),
		validation.Field(&r.Location, validation.Required),
		validation.Field(&r.StartDate, validation.Required),
		validation.Field(&r.EndDate, validation.Required),
		validation.Field(&r.Deadline, validation.Required),
		validation.Field(&r.SportID, validation.Required, is.UUID),
		validation.Field(&r.Rules, validation.Required),
		validation.Field(&r.ProposalLink, is.URL),
		validation.Field(&r.Open, validation.Required),
		validation.Field(&r.Quota, validation.Required),
	)
}

type StatusReq struct {
	Status *bool `json:"status"`
}

type response struct {
	ID            uuid.UUID                        `json:"id"`
	UserID        uuid.UUID                        `json:"user_id"`
	UserName      string                           `json:"user_name"`
	UserImage     string                           `json:"user_image"`
	Type          db.EventType                     `json:"type"`
	Name          string                           `json:"name"`
	Description   string                           `json:"description"`
	PrizePool     string                           `json:"prize_pool"`
	Location      string                           `json:"location"`
	Province      string                           `json:"province"`
	City          string                           `json:"city"`
	Thumbnail     string                           `json:"thumbnail"`
	StartDate     string                           `json:"start_date"`
	EndDate       string                           `json:"end_date"`
	Deadline      string                           `json:"deadline"`
	SportID       uuid.UUID                        `json:"sport_id"`
	SportName     string                           `json:"sport_name"`
	Rules         string                           `json:"rules"`
	ProposalLink  string                           `json:"proposal_link"`
	Status        bool                             `json:"status"`
	Quota         int32                            `json:"quota"`
	Open          string                           `json:"open"`
	Remark        string                           `json:"remark"`
	ClassEvents   []db.ClassEventFetchByEventIDRow `json:"class_events"`
	Participants  int64                            `json:"participants"`
	Privilege     db.EventPrivilegeFetchOneRow     `json:"user_privilege"`
	EventTurnLock bool                             `json:"event_turn_lock"`
}

type ResponseFetchOne struct {
	ID               uuid.UUID                      `json:"id"`
	UserID           uuid.UUID                      `json:"user_id"`
	UserName         string                         `json:"user_name"`
	UserImage        string                         `json:"user_image"`
	Type             db.EventType                   `json:"type"`
	Name             string                         `json:"name"`
	Description      string                         `json:"description"`
	PrizePool        string                         `json:"prize_pool"`
	Location         string                         `json:"location"`
	Province         string                         `json:"province"`
	City             string                         `json:"city"`
	Thumbnail        string                         `json:"thumbnail"`
	StartDate        string                         `json:"start_date"`
	EndDate          string                         `json:"end_date"`
	Deadline         time.Time                      `json:"deadline"`
	SportID          uuid.UUID                      `json:"sport_id"`
	SportName        string                         `json:"sport_name"`
	Rules            string                         `json:"rules"`
	ProposalLink     string                         `json:"proposal_link"`
	Status           bool                           `json:"status"`
	Quota            int32                          `json:"quota"`
	Open             string                         `json:"open"`
	Remark           string                         `json:"remark"`
	ClassEvents      []classEventSummaryRow         `json:"class_events"`
	Participants     int64                          `json:"participants"`
	Privilege        db.EventPrivilegeFetchOneRow   `json:"user_privilege"`
	EventTurnLock    bool                           `json:"event_turn_lock"`
	GeneralChampions []db.RankClubFetchByEventIDRow `json:"general_champions"`
}

type classEventSummaryRow struct {
	db.ClassEventFetchByEventIDRow
	Summary []RankFetchByClassEventIDRow `json:"summary"`
}

type FetchInfiniteQueryParams struct {
	SportID  string        `query:"sport_id"`
	Name     string        `query:"name"`
	Category db.SportType  `query:"category"`
	Order    int64         `query:"order"`
	Remark   db.RemarkType `query:"remark"`
}

type ResponseInfinite struct {
	Message   string                `json:"message"`
	Data      []ResponseInfiniteRow `json:"data"`
	TotalItem int64                 `json:"total_item"`
}

type ResponseInfiniteRow struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"user_id"`
	UserImage    string       `json:"user_image"`
	UserName     string       `json:"user_name"`
	Type         db.EventType `json:"type"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	PrizePool    string       `json:"prize_pool"`
	Location     string       `json:"location"`
	Province     string       `json:"province"`
	City         string       `json:"city"`
	Thumbnail    string       `json:"thumbnail"`
	StartDate    string       `json:"start_date"`
	EndDate      string       `json:"end_date"`
	Deadline     string       `json:"deadline"`
	SportID      uuid.UUID    `json:"sport_id"`
	SportName    string       `json:"sport_name"`
	Quota        int32        `json:"quota"`
	Order        int64        `json:"order"`
	Open         string       `json:"open"`
	Remark       string       `json:"remark"`
	Participants int64        `json:"participants"`
}

type AssignRequest struct {
	Data    []AssignDataRequest `json:"data"`
	EventID string              `param:"event"`
}

func (a AssignRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Data, validation.Required),
		validation.Field(&a.EventID, validation.Required, is.UUID),
	)
}

type AssignDataRequest struct {
	UserID string       `json:"user_id"`
	Role   db.EventRole `json:"role"`
}

func (a AssignDataRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Role, validation.Required),
		validation.Field(&a.UserID, validation.Required, is.UUID),
	)
}

type UpdateCommitteeParams struct {
	EventID     uuid.UUID    `param:"event"`
	CommitteeID uuid.UUID    `param:"committee"`
	Role        db.EventRole `json:"role"`
}

type DeleteCommitteeParams struct {
	EventID     uuid.UUID `param:"event"`
	CommitteeID uuid.UUID `param:"committee"`
}

type UpdateRemarkParams struct {
	EventID uuid.UUID     `param:"event"`
	Remark  db.RemarkType `json:"remark"`
}

func (u UpdateRemarkParams) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Remark, validation.Required),
	)
}
