package club

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository) Usecase {
	return Usecase{repository: repository, rawRepository: rawRepository}
}

func (u Usecase) store(ctx context.Context, req request, userID string) (uuid.UUID, error) {
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error in start db transaction")
	}
	txQuery := u.repository.WithTx(tx)
	clubID, err := txQuery.ClubCreate(ctx, db.ClubCreateParams{
		Name:      req.Name,
		Logo:      req.Logo,
		Phone:     req.Phone,
		ShortName: req.ShortName,
		UserID:    uuid.MustParse(userID),
	})
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return uuid.UUID{}, fmt.Errorf("error in rollback transaction %s", err.Error())
		}
		return uuid.UUID{}, fmt.Errorf("error in db transaction create club : %s", err.Error())
	}
	for _, sportData := range req.Sports {
		if err = txQuery.ClubAttachSport(ctx, db.ClubAttachSportParams{
			ClubID:  clubID,
			SportID: uuid.MustParse(sportData.SportID),
		}); err != nil {
			if err = tx.Rollback(); err != nil {
				return uuid.UUID{}, fmt.Errorf("error in rollback transaction %s", err.Error())
			}
			return uuid.UUID{}, fmt.Errorf("error in db transaction attach sport to club : %s", err.Error())
		}
	}
	if err = txQuery.ClubParticipantCreate(ctx, db.ClubParticipantCreateParams{
		ClubID: clubID,
		UserID: uuid.MustParse(userID),
	}); err != nil {
		if err = tx.Rollback(); err != nil {
			return uuid.UUID{}, fmt.Errorf("error in rollback transaction %s", err.Error())
		}
		return uuid.UUID{}, fmt.Errorf("error in db transaction add participant to club : %s", err.Error())
	}
	if err = tx.Commit(); err != nil {
		return uuid.UUID{}, fmt.Errorf("error in commit transaction")
	}
	return clubID, nil
}

func (u Usecase) update(ctx context.Context, req request, userID string, clubID uuid.UUID) error {
	ID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     clubID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check club : %s", err.Error())
	}
	return u.repository.ClubUpdate(ctx, db.ClubUpdateParams{
		Name:      req.Name,
		Logo:      req.Logo,
		Phone:     req.Phone,
		ShortName: req.ShortName,
		ID:        ID,
	})
}

func (u Usecase) delete(ctx context.Context, userID string, clubID uuid.UUID) error {
	ID, err := u.repository.ClubCheckOne(ctx, db.ClubCheckOneParams{
		ID:     clubID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check club : %s", err.Error())
	}
	return u.repository.ClubDelete(ctx, ID)
}
