package bracket

//go:generate mockgen -source=./usecase.go -destination=../../mocks/bracket/usecase_mock.go

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Usecase interface {
	Store(ctx context.Context, arg GenerateParams) (statusCode int, err error)
	FetchOne(ctx context.Context, arg GenerateParams) (FetchOneResponse, error)
	RoundDown(ctx context.Context, arg GenerateParams) (statusCode int, response RoundDownResponse, err error)
	OrderLock(ctx context.Context, arg UpdateLockParams) (statusCode int, err error)
	CancelBracket(ctx context.Context, arg UpdateGenerateParams) (statusCode int, err error)
	UpdateSingleLock(ctx context.Context, arg UpdateSingleLockParams) (statusCode int, err error)
	EventTurnLock(ctx context.Context, eventID uuid.UUID) (statusCode int, err error)
}

type usecase struct {
	repository    *db.Queries
	rawRepository Repository
	r             *rand.Rand
	db            *sql.DB
}

func NewUsecase(repository *db.Queries, rawRepository Repository, r *rand.Rand, db *sql.DB) Usecase {
	return &usecase{
		repository:    repository,
		rawRepository: rawRepository,
		r:             r,
		db:            db,
	}
}

func (u *usecase) Store(ctx context.Context, arg GenerateParams) (statusCode int, err error) {
	classEvent, err := u.repository.ClassEventFetchOne(ctx, arg.ClassEventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("class event not found : %s", err.Error())
	}

	if classEvent.BracketGenerate {
		return http.StatusForbidden, fmt.Errorf("bracket for specific class event already generated")
	}

	registrations, err := u.repository.EventRegistrationFetchByClassEventID(ctx, classEvent.ID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("registration not found : %s", err.Error())
	}

	tx, err := u.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)

	//Match Index
	matchIndex := int(math.Ceil(math.Log(float64(len(registrations))) / math.Log(2)))
	if err = txQuery.ClassEventUpdateBracketGenerate(ctx, db.ClassEventUpdateBracketGenerateParams{
		BracketGenerate: true,
		MatchIndex:      int16(matchIndex),
		ID:              classEvent.ID,
	}); err != nil {
		if err = tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update class event : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update class event : %s", err.Error())
	}
	switch classEvent.MatchType {
	case "order":
		for _, registration := range registrations {
			if err = txQuery.OrderBracketCreate(ctx, db.OrderBracketCreateParams{
				EventID:             classEvent.EventID,
				ClassEventID:        arg.ClassEventID,
				EventRegistrationID: registration.ID,
				ClubID:              registration.ClubID,
			}); err != nil {
				if err = tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx generate bracket order : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error in generate bracket order : %s", err.Error())
			}
		}
	case "single":
		count := len(registrations)
		arg := storeBracketParams{
			Ctx:               ctx,
			Tx:                tx,
			TxQuery:           txQuery,
			ClassEvent:        classEvent,
			RegistrationCount: count,
			MatchIndex:        matchIndex,
		}
		if count < 2 {
			statusCode, err = u.storeBracketDirectWinner(arg)
			if err != nil {
				return statusCode, err
			}
			break
		}
		statusCode, err = u.storeBracket(arg)
		if err != nil {
			return statusCode, err
		}

		//Store Next match id
		for i := 1; i < matchIndex; i++ {
			nextEventBrackets, err := txQuery.EventBracketFetchByMatchIndexAndClassID(ctx, db.EventBracketFetchByMatchIndexAndClassIDParams{
				MatchIndex:   int16(i),
				ClassEventID: classEvent.ID,
			})
			if err != nil {
				if err = tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx fetch event brackets : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error in fetch event brackets : %s", err.Error())
			}
			for _, eventBracket := range nextEventBrackets {
				nextMatchOrder := eventBracket.MatchOrder * 2
				if err = txQuery.EventBracketUpdateNextMatch(ctx, db.EventBracketUpdateNextMatchParams{
					NextMatchID:  eventBracket.ID,
					ClassEventID: classEvent.ID,
					MatchIndex:   int16(i) + 1,
					MatchOrder:   nextMatchOrder,
				}); err != nil {
					if err = tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update event next match : %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error in update event next match : %s", err.Error())
				}
				if err = txQuery.EventBracketUpdateNextMatch(ctx, db.EventBracketUpdateNextMatchParams{
					NextMatchID:  eventBracket.ID,
					ClassEventID: classEvent.ID,
					MatchIndex:   int16(i) + 1,
					MatchOrder:   nextMatchOrder - 1,
				}); err != nil {
					if err := tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update event next match : %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error in update event next match : %s", err.Error())
				}
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return http.StatusCreated, nil
}

func (u *usecase) storeBracketDirectWinner(arg storeBracketParams) (statusCode int, err error) {
	bracketID, err := arg.TxQuery.EventBracketCreate(arg.Ctx, db.EventBracketCreateParams{
		EventID:      arg.ClassEvent.EventID,
		ClassEventID: arg.ClassEvent.ID,
		EventTurn:    0,
		MatchIndex:   1,
		MatchOrder:   1,
		NextMatchID:  uuid.UUID{},
		IsActive:     0,
		Status:       db.BracketTypeBattle,
	})
	if err != nil {
		if err := arg.Tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create event bracket : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in create event bracket : %s", err.Error())
	}
	if err := arg.TxQuery.BracketParticipantCreate(arg.Ctx, db.BracketParticipantCreateParams{
		EventRegistrationID: uuid.UUID{},
		EventBracketID:      bracketID,
		Type:                db.ParticipantTypeHome,
		IsBye:               false,
	}); err != nil {
		if err := arg.Tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create bracket participant : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in create bracket participant : %s", err.Error())
	}
	if err := arg.TxQuery.BracketParticipantCreate(arg.Ctx, db.BracketParticipantCreateParams{
		EventRegistrationID: uuid.UUID{},
		EventBracketID:      bracketID,
		Type:                db.ParticipantTypeAway,
		IsBye:               false,
	}); err != nil {
		if err := arg.Tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create bracket participant : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in create bracket participant : %s", err.Error())
	}
	return http.StatusCreated, nil
}

func (u *usecase) storeBracket(arg storeBracketParams) (statusCode int, err error) {
	for index := 1; index <= arg.MatchIndex; index++ {
		//In one index how many bracket we must generate
		bracketIndex := int(math.Pow(2, float64(index)))
		if bracketIndex >= arg.RegistrationCount {
			bracketIndex = arg.RegistrationCount
		}
		matches := u.generateBracket(bracketIndex, index)

		var isActive int16
		if index == arg.MatchIndex {
			isActive = 1
		}

		//Store Match from Final to the latest and fill the id of next bracket
		for i := 1; i <= len(matches); i++ {
			bracketID, err := arg.TxQuery.EventBracketCreate(arg.Ctx, db.EventBracketCreateParams{
				EventID:      arg.ClassEvent.EventID,
				ClassEventID: arg.ClassEvent.ID,
				EventTurn:    0,
				MatchIndex:   int16(index),
				MatchOrder:   int16(i),
				NextMatchID:  uuid.UUID{},
				IsActive:     isActive,
				Status:       db.BracketTypeBattle,
			})
			if err != nil {
				if err := arg.Tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create event bracket : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error in create event bracket : %s", err.Error())
			}
			var home, away bool
			if matches[i-1][0] == 0 {
				home = true
			}
			if matches[i-1][1] == 0 {
				away = true
			}

			if home || away {
				if err := arg.TxQuery.EventBracketUpdateStatus(arg.Ctx, db.EventBracketUpdateStatusParams{
					Status:   db.BracketTypeBye,
					IsActive: 0,
					ID:       bracketID,
				}); err != nil {
					if err := arg.Tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update event bracket status : %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error in update event bracket status : %s", err.Error())
				}
			}

			if err := arg.TxQuery.BracketParticipantCreate(arg.Ctx, db.BracketParticipantCreateParams{
				EventRegistrationID: uuid.UUID{},
				EventBracketID:      bracketID,
				Type:                db.ParticipantTypeHome,
				IsBye:               home,
			}); err != nil {
				if err := arg.Tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create bracket participant : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error in create bracket participant : %s", err.Error())
			}
			if err = arg.TxQuery.BracketParticipantCreate(arg.Ctx, db.BracketParticipantCreateParams{
				EventRegistrationID: uuid.UUID{},
				EventBracketID:      bracketID,
				Type:                db.ParticipantTypeAway,
				IsBye:               away,
			}); err != nil {
				if err = arg.Tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create bracket participant : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error in create bracket participant : %s", err.Error())
			}
		}
	}
	return http.StatusCreated, nil
}

// References : https://stackoverflow.com/questions/5770990/sorting-tournament-seeds/45572051#45572051
func (u *usecase) generateBracket(registrationCount int, matchIndex int) [][]int {
	//bracketSize := math.Pow(2, matchIndex)
	//requiredByes := int(bracketSize) - registrationCount

	//Per match 1 vs 1
	matches := [][]int{{1, 2}}

	for round := 1; round < matchIndex; round++ {
		var roundMatches [][]int
		sum := math.Pow(2, float64(round+1)) + 1
		for _, match := range matches {
			home := u.changeIntoBye(match[0], registrationCount)
			away := u.changeIntoBye(int(sum)-match[0], registrationCount)
			roundMatches = append(roundMatches, []int{home, away})
			home = u.changeIntoBye(int(sum)-match[1], registrationCount)
			away = u.changeIntoBye(match[1], registrationCount)
			roundMatches = append(roundMatches, []int{home, away})
		}
		matches = roundMatches
	}
	return matches
}

func (u *usecase) changeIntoBye(seed int, participantCount int) int {
	if seed <= participantCount {
		return seed
	}
	return 0
}

func (u *usecase) FetchOne(ctx context.Context, arg GenerateParams) (FetchOneResponse, error) {
	classEvent, err := u.repository.ClassEventFetchOne(ctx, arg.ClassEventID)
	if err != nil {
		return FetchOneResponse{}, fmt.Errorf("class event not found : %s", err.Error())
	}
	var result interface{}
	switch classEvent.MatchType {
	case db.MatchTypeOrder:
		brackets, err := u.rawRepository.OrderBracketFetchByClassEventID(ctx, arg.ClassEventID)
		if err != nil {
			return FetchOneResponse{}, fmt.Errorf("brackets with specific classs event not found : %s", err.Error())
		}
		results := []fetchOneOrderResponse{}
		for _, bracket := range brackets {
			score, _ := u.repository.OrderScoreFetchOneByBracketID(ctx, bracket.ID)
			results = append(results, fetchOneOrderResponse{
				OrderBracketFetchByClassEventIDRow: bracket,
				Scores:                             score,
			})
		}
		result = results
	case db.MatchTypeSingle:
		response := []matchIndexData{}
		for index := classEvent.MatchIndex; index >= 1; index-- {
			indexTitle := u.setIndexTitle(index)
			matchIndexResult := matchIndexData{
				Title: indexTitle,
				Seeds: nil,
			}
			brackets, err := u.repository.EventBracketFetchByClassEventID(ctx, db.EventBracketFetchByClassEventIDParams{
				ClassEventID: classEvent.ID,
				MatchIndex:   index,
			})
			if err != nil {
				return FetchOneResponse{}, fmt.Errorf("error in fetch brackets : %s", err.Error())
			}
			for _, bracket := range brackets {
				var isScore bool
				if bracket.HomeTotal != 0 && bracket.AwayTotal != 0 {
					isScore = true
				}
				teams, err := u.rawRepository.BracketParticipantFetchByEventBracketID(ctx, bracket.ID)
				if err != nil {
					return FetchOneResponse{}, fmt.Errorf("error in fetch bracket participants : %s", err.Error())
				}
				var participants []bracketParticipantResponse
				for _, team := range teams {
					var participant bracketParticipantResponse
					participant.BracketParticipantFetchByEventBracketIDRow = team
					club, _ := u.repository.ClubFetchOne(ctx, team.ClubID)
					participant.ClubLogo = club.Logo
					switch team.Type {
					case db.ParticipantTypeHome:
						score, _ := u.repository.EventScoreFetchHomeByBracketID(ctx, bracket.ID)
						participant.Score = bracketParticipantScoreParams{
							Round1: score.HomeRound1,
							Round2: score.HomeRound1,
							Round3: score.HomeRound3,
							Extra:  score.HomeExtra,
							Total:  score.HomeTotal,
						}
						participants = append(participants, participant)
					case db.ParticipantTypeAway:
						score, _ := u.repository.EventScoreFetchAwayByBracketID(ctx, bracket.ID)
						participant.Score = bracketParticipantScoreParams{
							Round1: score.AwayRound1,
							Round2: score.AwayRound1,
							Round3: score.AwayRound3,
							Extra:  score.AwayExtra,
							Total:  score.AwayTotal,
						}
						participants = append(participants, participant)
					}
				}
				matchIndexResult.Seeds = append(matchIndexResult.Seeds, seedData{
					ID:         bracket.ID,
					EventTurn:  bracket.EventTurn,
					MatchOrder: bracket.MatchOrder,
					IsActive:   bracket.IsActive,
					IsScore:    isScore,
					Teams:      participants,
				})
			}
			response = append(response, matchIndexResult)
		}
		result = response
	}

	summary, _ := u.rawRepository.RankFetchByClassEventID(ctx, arg.ClassEventID)

	return FetchOneResponse{
		Message:        "fetch one bracket for specific event class success",
		Data:           result,
		GenerateStatus: &classEvent.BracketGenerate,
		LockStatus:     &classEvent.BracketLock,
		MatchType:      classEvent.MatchType,
		LockScore:      &classEvent.ScoreLock,
		Summary:        summary,
	}, nil
}

// reference = https://en.wikipedia.org/wiki/Single-elimination_tournament
func (u *usecase) setIndexTitle(matchIndex int16) string {
	var indexTitle string
	switch matchIndex {
	case 1:
		indexTitle = "Final"
	case 2:
		indexTitle = "Semifinals"
	case 3:
		indexTitle = "Quarterfinals"
	case 4:
		indexTitle = "Round of 16"
	case 5:
		indexTitle = "Round of 32"
	case 6:
		indexTitle = "Round of 64"
	case 7:
		indexTitle = "Round of 128"
	}
	return indexTitle
}

func (u *usecase) RoundDown(ctx context.Context, arg GenerateParams) (statusCode int, response RoundDownResponse, err error) {
	classEvent, err := u.repository.ClassEventFetchOne(ctx, arg.ClassEventID)
	if err != nil {
		return http.StatusNotFound, RoundDownResponse{}, fmt.Errorf("class event not found : %s", err.Error())
	}
	if classEvent.BracketLock {
		return http.StatusForbidden, RoundDownResponse{}, fmt.Errorf("bracket already been round down")
	}

	switch classEvent.MatchType {
	case db.MatchTypeOrder:
		brackets, err := u.rawRepository.OrderBracketFetchByClassEventID(ctx, arg.ClassEventID)
		if err != nil {
			return http.StatusNotFound, RoundDownResponse{}, fmt.Errorf("bracket not found : %s", err.Error())
		}
		u.r.Seed(time.Now().UnixNano())
		u.r.Shuffle(len(brackets), func(i, j int) {
			brackets[i], brackets[j] = brackets[j], brackets[i]
		})
		var result []orderRoundDownResponse
		for i, bracket := range brackets {
			result = append(result, orderRoundDownResponse{
				OrderBracketFetchByClassEventIDRow: bracket,
				Iteration:                          int16(i) + 1,
			})
		}
		response.Data = result
	case db.MatchTypeSingle:
		statusCode, result, err := u.roundDownBracketSingle(ctx, classEvent)
		if err != nil {
			return statusCode, RoundDownResponse{}, err
		}
		response.Data = result
	}
	return http.StatusOK, RoundDownResponse{
		Message:   "round down bracket success",
		Data:      response.Data,
		MatchType: classEvent.MatchType,
	}, nil
}

func (u *usecase) roundDownBracketSingle(ctx context.Context, classEvent db.ClassEventFetchOneRow) (statusCode int, result []matchIndexData, err error) {
	registrations, err := u.repository.EventRegistrationFetchByClassEventID(ctx, classEvent.ID)
	if err != nil {
		return http.StatusNotFound, nil, fmt.Errorf("registration data not found : %s", err.Error())
	}
	matches := u.generateBracket(len(registrations), int(classEvent.MatchIndex))
	u.r.Seed(time.Now().UnixNano())
	u.r.Shuffle(len(registrations), func(i, j int) {
		registrations[i], registrations[j] = registrations[j], registrations[i]
	})
	var randomTeams []BracketParticipantFetchByEventBracketIDRow
	var indexes []int
	min := 0
	max := len(registrations)
	for _, match := range matches {
		u.r.Seed(time.Now().UnixNano())
		index := rand.Intn(max-min) + min
		if match[0] != 0 {
			for {
				_, found := Find(indexes, index)
				if index == len(registrations) {
					break
				}
				if !found {
					break
				}
				index = rand.Intn(max-min) + min
			}
			registration := registrations[index]
			eventParticipants, err := u.repository.EventParticipantFetchNameByRegistrationID(ctx, registration.ID)
			if err != nil {
				return http.StatusNotFound, nil, fmt.Errorf("event participants not found : %s", err.Error())
			}
			randomTeams = append(randomTeams, BracketParticipantFetchByEventBracketIDRow{
				ID:     0,
				ClubID: registration.ClubID,
				ClubName: sql.NullString{
					String: registration.ClubName,
					Valid:  true,
				},
				Type:                db.ParticipantTypeHome,
				Participants:        eventParticipants,
				IsBye:               false,
				EventRegistrationID: registration.ID,
			})
			indexes = append(indexes, index)
		} else {
			randomTeams = append(randomTeams, BracketParticipantFetchByEventBracketIDRow{
				ID:     0,
				ClubID: uuid.UUID{},
				ClubName: sql.NullString{
					String: "",
					Valid:  false,
				},
				Type:                db.ParticipantTypeHome,
				Participants:        nil,
				IsBye:               true,
				EventRegistrationID: uuid.UUID{},
			})
		}
		if match[1] != 0 {
			for {
				_, found := Find(indexes, index)
				if !found {
					break
				}
				index = rand.Intn(max-min) + min
			}
			registration := registrations[index]
			eventParticipants, err := u.repository.EventParticipantFetchNameByRegistrationID(ctx, registration.ID)
			if err != nil {
				return http.StatusNotFound, nil, fmt.Errorf("event participants not found : %s", err.Error())
			}
			randomTeams = append(randomTeams, BracketParticipantFetchByEventBracketIDRow{
				ID:     0,
				ClubID: registration.ClubID,
				ClubName: sql.NullString{
					String: registration.ClubName,
					Valid:  true,
				},
				Type:                db.ParticipantTypeAway,
				Participants:        eventParticipants,
				IsBye:               false,
				EventRegistrationID: registration.ID,
			})
			indexes = append(indexes, index)
		} else {
			randomTeams = append(randomTeams, BracketParticipantFetchByEventBracketIDRow{
				ID:     0,
				ClubID: uuid.UUID{},
				ClubName: sql.NullString{
					String: "",
					Valid:  false,
				},
				Type:                db.ParticipantTypeAway,
				Participants:        nil,
				IsBye:               true,
				EventRegistrationID: uuid.UUID{},
			})
		}
	}
	result = []matchIndexData{}
	for index := classEvent.MatchIndex; index >= 1; index-- {
		indexTitle := u.setIndexTitle(index)
		matchIndexResult := matchIndexData{
			Title: indexTitle,
			Seeds: nil,
		}
		brackets, err := u.repository.EventBracketFetchByClassEventID(ctx, db.EventBracketFetchByClassEventIDParams{
			ClassEventID: classEvent.ID,
			MatchIndex:   index,
		})
		if err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("error in fetch brackets : %s", err.Error())
		}
		for _, bracket := range brackets {
			var teams []bracketParticipantResponse
			if index == classEvent.MatchIndex {
				for i := 1; i <= 2; i++ {
					teamIndex := (bracket.MatchOrder * 2) - 2
					switch i {
					case 1:
						teams = append(teams, bracketParticipantResponse{
							BracketParticipantFetchByEventBracketIDRow: randomTeams[teamIndex],
						})
					case 2:
						teams = append(teams, bracketParticipantResponse{
							BracketParticipantFetchByEventBracketIDRow: randomTeams[teamIndex+1],
						})
					}
				}
			} else {
				result, err := u.rawRepository.BracketParticipantFetchByEventBracketID(ctx, bracket.ID)
				if err != nil {
					return http.StatusInternalServerError, nil, fmt.Errorf("error inf etch bracket participants : %s", err.Error())
				}
				for _, r := range result {
					teams = append(teams, bracketParticipantResponse{
						BracketParticipantFetchByEventBracketIDRow: r,
					})
				}
			}
			matchIndexResult.Seeds = append(matchIndexResult.Seeds, seedData{
				ID:         bracket.ID,
				EventTurn:  bracket.EventTurn,
				MatchOrder: bracket.MatchOrder,
				Teams:      teams,
				IsActive:   bracket.IsActive,
			})
		}
		result = append(result, matchIndexResult)
	}
	return http.StatusOK, result, nil
}

func (u *usecase) OrderLock(ctx context.Context, arg UpdateLockParams) (statusCode int, err error) {
	classEvent, err := u.repository.ClassEventFetchOne(ctx, arg.ClassEventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("class event not found : %s", err.Error())
	}
	if classEvent.BracketLock == *arg.Status {
		return http.StatusBadRequest, fmt.Errorf("class event is already with lock status : %t", *arg.Status)
	}
	if classEvent.MatchType != db.MatchTypeOrder {
		return http.StatusBadRequest, fmt.Errorf("can't order a class event with match type except order")
	}
	tx, err := u.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	if err := txQuery.ClassEventUpdateBracketLock(ctx, db.ClassEventUpdateBracketLockParams{
		BracketLock: *arg.Status,
		ID:          classEvent.ID,
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket lock status : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update bracket lock status : %s", err.Error())
	}
	if *arg.Status {
		for _, participant := range arg.Participants {
			if err := u.repository.OrderBracketUpdateOrderBy(ctx, db.OrderBracketUpdateOrderByParams{
				OrderBy:             participant.Iteration,
				EventRegistrationID: uuid.MustParse(participant.EventRegistrationID),
			}); err != nil {
				if err := tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket lock status : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error in update bracket lock status : %s", err.Error())
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return http.StatusOK, nil
}

func (u *usecase) CancelBracket(ctx context.Context, arg UpdateGenerateParams) (statusCode int, err error) {
	classEvent, err := u.repository.ClassEventFetchOne(ctx, arg.ClassEventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("class event not found : %s", err.Error())
	}
	if !classEvent.BracketGenerate {
		return http.StatusBadRequest, fmt.Errorf("bracket already canceled")
	}
	tx, err := u.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	if err := txQuery.ClassEventUpdateBracketGenerate(ctx, db.ClassEventUpdateBracketGenerateParams{
		BracketGenerate: false,
		ID:              classEvent.ID,
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update class event : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update class event : %s", err.Error())
	}
	if err := txQuery.OrderBracketDeleteByClassEventID(ctx, classEvent.ID); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx delete order bracket : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in delete order bracket : %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}
	return http.StatusOK, nil
}

func (u *usecase) UpdateSingleLock(ctx context.Context, arg UpdateSingleLockParams) (statusCode int, err error) {
	classEvent, err := u.repository.ClassEventFetchOne(ctx, arg.ClassEventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("class event not found : %s", err.Error())
	}
	if classEvent.BracketLock == *arg.Status {
		return http.StatusBadRequest, fmt.Errorf("class event is already with lock status : %t", *arg.Status)
	}
	if classEvent.MatchType != db.MatchTypeSingle {
		return http.StatusBadRequest, fmt.Errorf("can't order a class event with match type except single")
	}

	bracketSize := int(math.Pow(2, float64(classEvent.MatchIndex))) / 2
	if len(arg.Data.Seeds) != bracketSize {
		return http.StatusBadRequest, fmt.Errorf("request data seed with the class event bracket is not the same")
	}

	tx, err := u.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)

	if err = txQuery.ClassEventUpdateBracketLock(ctx, db.ClassEventUpdateBracketLockParams{
		BracketLock: *arg.Status,
		ID:          classEvent.ID,
	}); err != nil {
		if err = tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket lock status : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update bracket lock status : %s", err.Error())
	}

	if arg.Status == nil {
		*arg.Status = false
	}

	if *arg.Status {
		for _, seed := range arg.Data.Seeds {
			for _, team := range seed.Teams {
				bracketParticipantID, err := txQuery.BracketParticipantCheckOneByEventBracketIDAndType(ctx, db.BracketParticipantCheckOneByEventBracketIDAndTypeParams{
					EventBracketID: seed.ID,
					Type:           team.Type,
				})
				if err != nil {
					return http.StatusNotFound, fmt.Errorf("bracket participant not found : %s", err.Error())
				}
				if err = txQuery.BracketParticipantUpdate(ctx, db.BracketParticipantUpdateParams{
					EventRegistrationID: team.EventRegistrationID,
					ID:                  bracketParticipantID,
				}); err != nil {
					if err := tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket participant : %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error in update bracket participant : %s", err.Error())
				}
			}
		}
	}

	//Store For Bracket Bye
	brackets, _ := txQuery.EventBracketFetchByClassEventStatusByeAndMatchIndex(ctx, db.EventBracketFetchByClassEventStatusByeAndMatchIndexParams{
		ClassEventID: arg.ClassEventID,
		MatchIndex:   classEvent.MatchIndex,
	})
	if len(brackets) > 0 {
		for _, bracket := range brackets {
			participants, err := txQuery.BracketParticipantFetchByEventBracketID(ctx, bracket.ID)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx bracket participants not found : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("bracket participants not found : %s", err.Error())
			}

			for _, participant := range participants {
				if participant.EventRegistrationID.String() != "00000000-0000-0000-0000-000000000000" {
					var nextParticipantType db.ParticipantType
					if bracket.MatchOrder%2 == 1 {
						nextParticipantType = db.ParticipantTypeHome
					} else {
						nextParticipantType = db.ParticipantTypeAway
					}

					if err := txQuery.BracketParticipantUpdateByParticipantType(ctx, db.BracketParticipantUpdateByParticipantTypeParams{
						EventRegistrationID: participant.EventRegistrationID,
						EventBracketID:      bracket.NextMatchID,
						Type:                nextParticipantType,
					}); err != nil {
						if err := tx.Rollback(); err != nil {
							return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket participants update : %s", err.Error())
						}
						return http.StatusInternalServerError, fmt.Errorf("error in update bracket participants : %s", err.Error())
					}
				}

				//Next Bracket
				if err := txQuery.EventBracketUpdateStatus(ctx, db.EventBracketUpdateStatusParams{
					Status:   db.BracketTypeBattle,
					IsActive: 1,
					ID:       bracket.NextMatchID,
				}); err != nil {
					if err := tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update bracket status : %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error in update bracket status : %s", err.Error())
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}
	return http.StatusOK, nil
}

func (u *usecase) EventTurnLock(ctx context.Context, eventID uuid.UUID) (statusCode int, err error) {
	event, err := u.repository.EventCheckOne(ctx, eventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("event not found : %s", err.Error())
	}

	if event.IsGenerate {
		return http.StatusBadRequest, fmt.Errorf("event already generate event turn")
	}

	classEvent, err := u.repository.ClassEventFetchSingleByEventIDAndLockStatus(ctx, eventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("error in fetch class event : %s", err.Error())
	}
	if len(classEvent) != 0 {
		return http.StatusForbidden, fmt.Errorf("to generate event turn, must generate all bracket")
	}
	classEventSingles, err := u.repository.ClassEventFetchAndGroupByEventID(ctx, eventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("fetch class event error : %s", err.Error())
	}

	lastMatchIndex, err := u.repository.ClassEventFetchLastMatchIndexByEventID(ctx, eventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("fetch last match index error : %s", err.Error())
	}

	tx, err := u.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	var eventTurn int16
	for index := lastMatchIndex; index >= 1; index-- {
		for _, classEventSingle := range classEventSingles {
			eventBrackets, err := txQuery.EventBracketFetchByClassEventStatusBattleAndMatchIndex(ctx, db.EventBracketFetchByClassEventStatusBattleAndMatchIndexParams{
				ClassEventID: classEventSingle.ID,
				MatchIndex:   int16(index),
			})
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx fetch event bracket : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error fetch event bracket : %s", err.Error())
			}
			for _, bracket := range eventBrackets {
				eventTurn += 1
				if err := txQuery.EventBracketUpdateEventTurnByID(ctx, db.EventBracketUpdateEventTurnByIDParams{
					EventTurn: eventTurn,
					ID:        bracket.ID,
				}); err != nil {
					if err := tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update event turn : %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error update event turn : %s", err.Error())
				}
			}
		}
	}

	if err := u.repository.EventUpdateIsGenerate(ctx, db.EventUpdateIsGenerateParams{
		IsGenerate: true,
		ID:         event.ID,
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update event : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update event : %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return http.StatusOK, nil
}

func Find(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
