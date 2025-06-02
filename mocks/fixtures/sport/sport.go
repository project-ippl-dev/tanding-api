package sportFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/sport"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

var (
	SportFetchAllResponse = tools.Pagination{
		TotalItem: 1,
		PageSize:  1,
		Page:      1,
		Data: []sport.FetchAllRow{
			{
				ID:          uuid.New(),
				Name:        "name",
				Description: "description",
				Type:        "type",
				Thumbnail:   "https://google.com",
			},
		},
	}
)
