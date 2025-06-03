package dbFixtures

import (
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

var (
	ClubParticipantFetchJoinApprovalRow = []db.ClubParticipantFetchJoinApprovalRow{
		{
			ID:        1,
			SportID:   uuid.New(),
			SportName: "sport name",
			Name:      "name",
		},
	}

	ClubParticipantFetchInviteApprovalRow = []db.ClubParticipantFetchInviteApprovalRow{
		{
			ID:        1,
			ClubID:    uuid.New(),
			SportID:   uuid.New(),
			SportName: "sport name",
			Name:      "name",
		},
	}
)
