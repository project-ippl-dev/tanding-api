package classFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	ClassFetchAllResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.ClassFetchAllRow{
			{
				ID:                     uuid.New(),
				Name:                   "class-name",
				ClassCompetitionRuleID: 1,
				ClassRuleName:          "class-rule-name",
				Type:                   db.ClassTypeDefault,
				SportID:                uuid.New(),
				MatchType:              db.MatchTypeSingle,
				SportName:              "sport-name",
			},
		},
	}

	ClassFetchBySportIDResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []db.ClassFetchBySportIDRow{
			{
				ID:                     uuid.New(),
				Name:                   "class-name",
				ClassCompetitionRuleID: 1,
				ClassRuleName:          "class-rule-name",
				Type:                   db.ClassTypeDefault,
				SportID:                uuid.New(),
				MatchType:              db.MatchTypeSingle,
				SportName:              "sport-name",
			},
		},
	}
)
