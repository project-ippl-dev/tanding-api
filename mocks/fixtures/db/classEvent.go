package dbFixtures

import (
	"github.com/google/uuid"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/project-ippl-dev/tanding-api/internal/db"
	eventFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/event"
)

type MockResponseClassEventFetchOneRowReq struct {
	MatchType         db.MatchType
	IsBracketGenerate bool
	IsBracketLock     bool
	IsScoreLock       bool
	RuleMale          int16
	RuleFemale        int16
	RuleTotal         int16
	MatchIndex        int16
}

func NewMockResponseClassEventFetchOneRow(mockReq MockResponseClassEventFetchOneRowReq) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "match_type", "bracket_generate", "bracket_lock", "score_lock", "price", "rule_male", "rule_female", "rule_total", "event_id", "match_index",
	}).AddRow(uuid.NewString(), string(mockReq.MatchType), mockReq.IsBracketGenerate, mockReq.IsBracketLock, mockReq.IsScoreLock, 1000, mockReq.RuleMale, mockReq.RuleFemale, mockReq.RuleTotal, eventFixtures.EventID.String(), mockReq.MatchIndex)

}

var (
	ClassEventFetchAllRows = []db.ClassEventFetchAllRow{
		{
			ID:        uuid.New(),
			ClassName: "class event name",
		},
	}
)
