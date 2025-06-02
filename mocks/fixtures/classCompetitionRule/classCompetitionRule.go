package classCompetitionRuleFixtures

import (
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	ClassCompetitionRuleFetchResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.ClassRuleFetchAllRow{
			{
				ID:     1,
				Name:   "class competition rule",
				Male:   1,
				Female: 0,
				Total:  1,
			},
		},
	}
)
