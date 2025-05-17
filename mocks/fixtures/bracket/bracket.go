package bracketFixtures

import (
	"github.com/google/uuid"
	eventRegistrationFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/eventRegistration"

	"github.com/project-ippl-dev/tanding-api/internal/bracket"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/utils/pointer"
)

var (
	ClassEventID              = uuid.New()
	OrderBracketID            = uuid.New()
	ClubID                    = uuid.New()
	ClubName                  = "ClubName"
	ClubLogo                  = "https://logo.com"
	ParticipantMatchTypeOrder = []string{"participant1"}

	OrderBracketFetchByClassEventIDRowTypeOrder = bracket.OrderBracketFetchByClassEventIDRow{
		ID:                  OrderBracketID,
		Rank:                1,
		Participants:        ParticipantMatchTypeOrder,
		ClubName:            ClubName,
		ClubLogo:            ClubLogo,
		EventRegistrationID: eventRegistrationFixtures.EventRegistrationID,
	}

	OrderScoreFetchOneByBracketIDRow = db.OrderScoreFetchOneByBracketIDRow{
		ID:    uuid.New(),
		Total: 100,
	}

	DataMatchTypeOrder = []bracket.FetchOneOrderResponse{
		{
			OrderBracketFetchByClassEventIDRow: OrderBracketFetchByClassEventIDRowTypeOrder,
			Scores:                             OrderScoreFetchOneByBracketIDRow,
		},
	}

	SummariesMatchTypeOrder = []bracket.RankFetchByClassEventIDRow{
		{
			ID:           uuid.New(),
			Rank:         1,
			Point:        10,
			ClubName:     ClubName,
			ClubLogo:     ClubLogo,
			Participants: ParticipantMatchTypeOrder,
		},
	}

	FetchOneResponseMatchTypeOrder = bracket.FetchOneResponse{
		Message:        "fetch one bracket for specific event class success",
		Data:           DataMatchTypeOrder,
		GenerateStatus: pointer.ConvertToPointer(true),
		LockStatus:     pointer.ConvertToPointer(true),
		MatchType:      db.MatchTypeOrder,
		LockScore:      pointer.ConvertToPointer(true),
		Summary:        SummariesMatchTypeOrder,
	}

	OrderRoundDownResponse = []bracket.OrderRoundDownResponse{
		{
			OrderBracketFetchByClassEventIDRow: OrderBracketFetchByClassEventIDRowTypeOrder,
			Iteration:                          1,
		},
	}

	RoundDownResponseMatchTypeOrder = bracket.RoundDownResponse{
		Message:   "round down bracket success",
		Data:      OrderRoundDownResponse,
		MatchType: db.MatchTypeOrder,
	}

	BracketMatchIndexData = bracket.MatchIndexData{
		Title: "Final",
		Seeds: []bracket.SeedData{
			{
				ID:         uuid.New(),
				EventTurn:  1,
				MatchOrder: 1,
				IsActive:   1,
				Teams: []bracket.BracketParticipantResponse{
					{
						BracketParticipantFetchByEventBracketIDRow: bracket.BracketParticipantFetchByEventBracketIDRow{
							ID:                  1,
							ClubID:              ClubID,
							Type:                db.ParticipantTypeHome,
							Participants:        []string{"participant1"},
							IsBye:               false,
							EventRegistrationID: eventRegistrationFixtures.EventRegistrationID,
						},
						ClubLogo: ClubLogo,
					},
					{
						BracketParticipantFetchByEventBracketIDRow: bracket.BracketParticipantFetchByEventBracketIDRow{
							ID:                  2,
							ClubID:              ClubID,
							Type:                db.ParticipantTypeAway,
							Participants:        []string{"participant2"},
							IsBye:               false,
							EventRegistrationID: eventRegistrationFixtures.EventRegistrationID,
						},
						ClubLogo: ClubLogo,
					},
				},
			},
		},
	}
)
