package accomplishmentFixtures

import (
	"github.com/project-ippl-dev/tanding-api/internal/accomplishment"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	AccomplishmentFetchAllResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []accomplishment.Response{
			{
				Title:       "title",
				Level:       "level",
				Ranking:     "ranking",
				Category:    "category",
				Sport:       "sport",
				Description: "description",
				FileURL:     "https://google.com",
			},
		},
	}
)
