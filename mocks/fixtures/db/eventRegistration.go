package dbFixtures

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	bracketFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/bracket"
	eventFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/event"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type MockResponseEventRegistrationFetchByClassEventIDReq struct {
	RegistrationIteration int
}

func NewMockResponseEventRegistrationFetchByClassEventID(req MockResponseEventRegistrationFetchByClassEventIDReq) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "event_id", "class_event_id", "club_id", "status", "club_name",
	})

	for i := 0; i < req.RegistrationIteration; i++ {
		rows.AddRow(uuid.NewString(), eventFixtures.EventID.String(), bracketFixtures.ClassEventID.String(), uuid.NewString(), string(db.EventRegistrationStatusApproved), fmt.Sprintf("club-%d", i))
	}

	return rows
}
