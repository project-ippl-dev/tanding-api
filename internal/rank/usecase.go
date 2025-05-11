package rank

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return Usecase{repository: repository, rawRepository: rawRepository}
}

func (u Usecase) summary(ctx context.Context, eventID uuid.UUID) (statusCode int, err error) {
	event, err := u.repository.EventFetchOne(ctx, eventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("event not found : %s", err.Error())
	}

	if event.Remark != db.RemarkTypeOngoing {
		return http.StatusForbidden, fmt.Errorf("already generate event summary")
	}

	classEvents, err := u.repository.ClassEventFetchByEventIDAndScoreLockTrue(ctx, event.ID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("class event not found : %s", err.Error())
	}

	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	for _, classEvent := range classEvents {
		var excludedRegistrationID []uuid.UUID
		switch classEvent.MatchType {
		case db.MatchTypeOrder:
			orderBrackets, err := txQuery.OrderBracketFetchByClassEventIDAndOrderByScore(ctx, classEvent.ID)
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error tx rollback order brackets not found : %s", err.Error())
				}
				return http.StatusNotFound, fmt.Errorf("order brackets not found : %s", err.Error())
			}
			for i := 0; i < 4; i++ {
				if i >= len(orderBrackets) {
					break
				}
				rank := int16(i + 1)
				point := u.setRankPoint(rank)
				if err := txQuery.RankCreate(ctx, db.RankCreateParams{
					ClubID:              orderBrackets[i].ClubID,
					EventID:             eventID,
					ClassEventID:        classEvent.ID,
					EventRegistrationID: orderBrackets[i].EventRegistrationID,
					SportID:             event.SportID,
					Rank:                rank,
					Point:               point,
				}); err != nil {
					if err := tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx rank create %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error in rank create : %s", err.Error())
				}
				if err := txQuery.OrderBracketUpdateRank(ctx, db.OrderBracketUpdateRankParams{
					Rank: rank,
					ID:   orderBrackets[i].ID,
				}); err != nil {
					if err := tx.Rollback(); err != nil {
						return http.StatusInternalServerError, fmt.Errorf("error in rollback tx order bracket update rank %s", err.Error())
					}
					return http.StatusInternalServerError, fmt.Errorf("error in order brakcet update rank : %s", err.Error())
				}
				statusCode, err = u.storeCertificateByRegistrationID(ctx, tx, txQuery, storeCertificateByRegistrationIDParams{
					EventRegistrationID: orderBrackets[i].EventRegistrationID,
					Rank:                rank,
					EventID:             event.ID,
					ClassName:           classEvent.ClassName,
					ClassEventID:        classEvent.ID,
				})
				if err != nil {
					return statusCode, err
				}

				excludedRegistrationID = append(excludedRegistrationID, orderBrackets[i].EventRegistrationID)
			}
			break
		case db.MatchTypeSingle:
			//Fetch Based on Match Index
			for i := 1; i <= 2; i++ {
				brackets, err := txQuery.EventBracketFetchByClassEventIDAndMatchIndex(ctx, db.EventBracketFetchByClassEventIDAndMatchIndexParams{
					ClassEventID: classEvent.ID,
					MatchIndex:   int16(i),
				})
				if err != nil {
					return http.StatusNotFound, fmt.Errorf("brackets not found : %s", err.Error())
				}
				if int16(i) > classEvent.MatchIndex {
					break
				}
				switch i {
				case 1:
					var homeRank, awayRank int16
					var homePoint, awayPoint int32
					if brackets[0].HomeTotal > brackets[0].AwayTotal {
						homeRank, awayRank = 1, 2
					} else {
						homeRank, awayRank = 2, 1
					}

					homePoint = u.setRankPoint(homeRank)
					awayPoint = u.setRankPoint(awayRank)

					bracketParticipants, err := txQuery.BracketParticipantFetchByEventBracketID(ctx, brackets[0].ID)
					if err != nil {
						return http.StatusNotFound, fmt.Errorf("bracket participant not found : %s", err.Error())
					}
					if bracketParticipants[0].EventRegistrationID.String() != "00000000-0000-0000-0000-000000000000" {
						if err := txQuery.RankCreate(ctx, db.RankCreateParams{
							ClubID:              bracketParticipants[0].ClubID,
							EventID:             eventID,
							ClassEventID:        classEvent.ID,
							EventRegistrationID: bracketParticipants[0].EventRegistrationID,
							SportID:             event.SportID,
							Rank:                homeRank,
							Point:               homePoint,
						}); err != nil {
							if err := tx.Rollback(); err != nil {
								return http.StatusInternalServerError, fmt.Errorf("error in rollback tx rank index 1 home create single : %s", err.Error())
							}
							return http.StatusInternalServerError, fmt.Errorf("error rank index 1 home create single : %s", err.Error())
						}
						statusCode, err = u.storeCertificateByRegistrationID(ctx, tx, txQuery, storeCertificateByRegistrationIDParams{
							EventRegistrationID: bracketParticipants[0].EventRegistrationID,
							Rank:                homeRank,
							EventID:             event.ID,
							ClassName:           classEvent.ClassName,
							ClassEventID:        classEvent.ID,
						})
						if err != nil {
							return statusCode, err
						}
						excludedRegistrationID = append(excludedRegistrationID, bracketParticipants[0].EventRegistrationID)
					}
					if bracketParticipants[1].EventRegistrationID.String() != "00000000-0000-0000-0000-000000000000" {
						if err := txQuery.RankCreate(ctx, db.RankCreateParams{
							ClubID:              bracketParticipants[1].ClubID,
							EventID:             eventID,
							ClassEventID:        classEvent.ID,
							EventRegistrationID: bracketParticipants[1].EventRegistrationID,
							SportID:             event.SportID,
							Rank:                awayRank,
							Point:               awayPoint,
						}); err != nil {
							if err := tx.Rollback(); err != nil {
								return http.StatusInternalServerError, fmt.Errorf("error in rollback tx rank index 1 away create single : %s", err.Error())
							}
							return http.StatusInternalServerError, fmt.Errorf("error rank index 1 away create single : %s", err.Error())
						}
						statusCode, err = u.storeCertificateByRegistrationID(ctx, tx, txQuery, storeCertificateByRegistrationIDParams{
							EventRegistrationID: bracketParticipants[1].EventRegistrationID,
							Rank:                awayRank,
							EventID:             event.ID,
							ClassName:           classEvent.ClassName,
							ClassEventID:        classEvent.ID,
						})
						if err != nil {
							return statusCode, err
						}
						excludedRegistrationID = append(excludedRegistrationID, bracketParticipants[1].EventRegistrationID)
					}
					break
				case 2:
					for j := 0; j < 2; j++ {
						if j+1 > len(brackets) {
							break
						}
						var rank int16 = 3
						point := u.setRankPoint(rank)
						var participantType db.ParticipantType
						if brackets[j].HomeTotal > brackets[j].AwayTotal {
							participantType = db.ParticipantTypeAway
						} else {
							participantType = db.ParticipantTypeHome
						}
						bracketParticipant, err := txQuery.BracketParticipantFetchByBracketIDAndType(ctx, db.BracketParticipantFetchByBracketIDAndTypeParams{
							EventBracketID: brackets[j].ID,
							Type:           participantType,
						})
						if err != nil {
							return http.StatusNotFound, fmt.Errorf("bracket participant not found : %s", err.Error())
						}
						if bracketParticipant.EventRegistrationID.String() != "00000000-0000-0000-0000-000000000000" {
							if err := txQuery.RankCreate(ctx, db.RankCreateParams{
								ClubID:              bracketParticipant.ClubID,
								EventID:             eventID,
								ClassEventID:        classEvent.ID,
								EventRegistrationID: bracketParticipant.EventRegistrationID,
								SportID:             event.SportID,
								Rank:                rank,
								Point:               point,
							}); err != nil {
								if err := tx.Rollback(); err != nil {
									return http.StatusInternalServerError, fmt.Errorf("error in rollback tx rank index 2 create single : %s", err.Error())
								}
							}
							statusCode, err = u.storeCertificateByRegistrationID(ctx, tx, txQuery, storeCertificateByRegistrationIDParams{
								EventRegistrationID: bracketParticipant.EventRegistrationID,
								Rank:                rank,
								EventID:             event.ID,
								ClassName:           classEvent.ClassName,
								ClassEventID:        classEvent.ID,
							})
							if err != nil {
								return statusCode, err
							}
							excludedRegistrationID = append(excludedRegistrationID, bracketParticipant.EventRegistrationID)
						}
					}
					break
				}
			}
			break
		}

		//Store Certificate For Participant
		statusCode, err := u.storeCertificateExcludeRegistrationID(ctx, tx, txQuery, storeCertificateExcludeRegistrationIDParams{
			EventRegistrationID: excludedRegistrationID,
			Rank:                0,
			EventID:             event.ID,
			ClassName:           classEvent.ClassName,
			ClassEventID:        classEvent.ID,
		})
		if err != nil {
			return statusCode, err
		}
	}

	statusCode, err = u.storeCertificateEventCommittee(ctx, tx, txQuery, eventID)
	if err != nil {
		return statusCode, err
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in first commit tx : %s", err.Error())
	}

	tx2, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start 2nd tx : %s", err.Error())
	}
	txQuery2 := u.repository.WithTx(tx2)

	statusCode, err = u.storeCertificateClubs(ctx, tx2, txQuery2, event.ID)
	if err != nil {
		return statusCode, err
	}

	if err := txQuery2.EventUpdateRemark(ctx, db.EventUpdateRemarkParams{
		Remark: db.RemarkTypeDone,
		ID:     event.ID,
	}); err != nil {
		if err := tx2.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update event remark : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update event remark")
	}

	if err := tx2.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in second commit tx : %s", err.Error())
	}

	return http.StatusCreated, nil
}

func (u Usecase) setRewardCertificateName(rank int16, className string) string {
	var rewardAs string
	switch rank {
	case 1:
		rewardAs = "Juara 1"
		break
	case 2:
		rewardAs = "Juara 2"
		break
	case 3, 4:
		rewardAs = "Juara 3"
		break
	default:
		rewardAs = "Peserta"
		break
	}
	return fmt.Sprintf("%s (%s)", rewardAs, className)
}

func (u *Usecase) setRewardCertificateCommitteeName(role db.EventRole) string {
	var rewardAs string
	switch role {
	case db.EventRoleOwner:
		rewardAs = "Event Organizer"
		break
	case db.EventRoleAdmin, db.EventRoleContributor:
		rewardAs = "Event Committee"
		break
	}
	return rewardAs
}

func (u Usecase) setRankPoint(rank int16) int32 {
	var point int
	switch rank {
	case 1:
		point = 10
		break
	case 2:
		point = 5
		break
	case 3, 4:
		point = 3
		break
	}
	return int32(point)
}

func (u *Usecase) storeCertificateByRegistrationID(ctx context.Context, tx *sql.Tx, txQuery *db.Queries, arg storeCertificateByRegistrationIDParams) (statusCode int, err error) {
	participants, err := txQuery.EventParticipantFetchByRegistrationID(ctx, arg.EventRegistrationID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx participants not found : %s", err.Error())
		}
		return http.StatusNotFound, fmt.Errorf("participants not found : %s", err.Error())
	}
	for _, participant := range participants {
		rewardAs := u.setRewardCertificateName(arg.Rank, arg.ClassName)
		if err := txQuery.CertificateCreate(ctx, db.CertificateCreateParams{
			ID:           uuid.New(),
			UserID:       participant.UserID,
			EventID:      arg.EventID,
			RewardAs:     rewardAs,
			ClassEventID: arg.ClassEventID,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				if err := tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create certificate : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error create certificate : %s", err.Error())
			}
		}
	}
	return http.StatusOK, nil
}

func (u Usecase) storeCertificateExcludeRegistrationID(ctx context.Context, tx *sql.Tx, txQuery *db.Queries, arg storeCertificateExcludeRegistrationIDParams) (statusCode int, err error) {
	participants, err := u.rawRepository.EventParticipantFetchExcludeRegistrationID(ctx, EventParticipantFetchExcludeRegistrationIDParams{
		EventRegistrationID: arg.EventRegistrationID,
		ClassEventID:        arg.ClassEventID,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in fetch excluded participant : %s", err.Error())
	}
	for _, participant := range participants {
		rewardAs := u.setRewardCertificateName(arg.Rank, arg.ClassName)
		if err := txQuery.CertificateCreate(ctx, db.CertificateCreateParams{
			ID:           uuid.New(),
			UserID:       participant.UserID,
			EventID:      arg.EventID,
			RewardAs:     rewardAs,
			ClassEventID: arg.ClassEventID,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				if err := tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create certificate : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error create certificate : %s", err.Error())
			}
		}
	}
	return http.StatusOK, nil
}

func (u Usecase) storeCertificateEventCommittee(ctx context.Context, tx *sql.Tx, txQuery *db.Queries, eventID uuid.UUID) (statusCode int, err error) {
	committees, err := txQuery.EventPrivilegeFetchByEventID(ctx, eventID)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx event committee not found : %s", err.Error())
		}
		return http.StatusNotFound, fmt.Errorf("event committee not found : %s", err.Error())
	}
	for _, committee := range committees {
		rewardAs := u.setRewardCertificateCommitteeName(committee.Role)
		if err := txQuery.CertificateCreate(ctx, db.CertificateCreateParams{
			ID:       uuid.New(),
			UserID:   committee.UserID,
			EventID:  eventID,
			RewardAs: rewardAs,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				if err := tx.Rollback(); err != nil {
					return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create certificate : %s", err.Error())
				}
				return http.StatusInternalServerError, fmt.Errorf("error create certificate : %s", err.Error())
			}
		}
	}
	return http.StatusOK, nil
}

func (u Usecase) storeCertificateClubs(ctx context.Context, tx *sql.Tx, txQuery *db.Queries, eventID uuid.UUID) (statusCode int, err error) {
	clubs, err := txQuery.RankClubFetchByEventID(ctx, eventID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in fetch clubs : %s", err.Error())
	}

	for i, club := range clubs {
		if i > 2 {
			break
		}
		rank := i + 1
		rewardAs := u.setCertificateClubName(rank)
		if err := txQuery.ClubCertificateCreate(ctx, db.ClubCertificateCreateParams{
			ID:       uuid.New(),
			ClubID:   club.ClubID,
			EventID:  eventID,
			RewardAs: rewardAs,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create club certificate : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in create club certificate : %s", err.Error())
		}
	}

	return http.StatusOK, nil
}

func (u Usecase) setCertificateClubName(rank int) string {
	var result string
	switch rank {
	case 1:
		result = "Juara Umum 1"
		break
	case 2:
		result = "Juara Umum 2"
		break
	case 3:
		result = "Juara Umum 3"
		break
	}
	return result
}

func (u Usecase) fetchOwnPoint(ctx context.Context, userID string) (int32, error) {
	return u.repository.RankFetchPointByUserID(ctx, uuid.MustParse(userID))
}

func (u Usecase) fetchByClubID(ctx context.Context, clubID uuid.UUID) (fetchByClubIDResponse, error) {
	total, _ := u.repository.RankFetchPointByClubID(ctx, clubID)
	participants, err := u.repository.RankFetchAllPointByClubID(ctx, clubID)
	if err != nil {
		return fetchByClubIDResponse{}, fmt.Errorf("error in fetch all point by club : %s", err.Error())
	}
	return fetchByClubIDResponse{
		TotalPoint:   total,
		Participants: participants,
	}, nil
}

func (u Usecase) rank(ctx context.Context, page, pageSize int32, arg rankParams) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)

	count, err := u.rawRepository.RankCountPowerList(ctx, arg.SportID)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count rank power list : %s", err.Error())
	}
	ranks, err := u.rawRepository.RankFetchPowerList(ctx, fetchPowerListParams{
		SportID: arg.SportID,
		Limit:   pageSize,
		Offset:  skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch rank power list : %s", err.Error())
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      ranks,
	}, nil
}

func (u Usecase) userRank(ctx context.Context, page, pageSize int32, arg rankParams) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)

	count, err := u.rawRepository.RankCountAllPointUser(ctx, arg.SportID)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count rank power list : %s", err.Error())
	}
	ranks, err := u.rawRepository.RankFetchAllPointUser(ctx, fetchAllPointUserParams{
		SportID: arg.SportID,
		Limit:   pageSize,
		Offset:  skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch rank power list : %s", err.Error())
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      ranks,
	}, nil
}
