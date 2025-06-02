package dbFixtures

import "github.com/project-ippl-dev/tanding-api/internal/db"

var (
	ClassRuleFetchOneRow = db.ClassRuleFetchOneRow{
		ID:     1,
		Name:   "class rule",
		Male:   1,
		Female: 0,
		Total:  1,
	}
)
