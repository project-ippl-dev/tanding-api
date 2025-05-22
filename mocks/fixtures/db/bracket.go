package dbFixtures

import (
	"github.com/google/uuid"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func NewMockResponseEventBracketCreate(bracketID uuid.UUID) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id"}).AddRow(bracketID.String())
}
