package classEvent

//go:generate mockgen -source=./usecase.go -destination=../../mocks/classEvent/usecase_mock.go

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/project-ippl-dev/tanding-api/internal/db"
)

type Usecase interface {
	Assign(ctx context.Context, userID string, req Request, eventID uuid.UUID) error
	Detach(ctx context.Context, userID string, eventID uuid.UUID, classID uuid.UUID) error
	FetchAll(ctx context.Context, userID string, eventID uuid.UUID) ([]db.ClassEventFetchAllRow, error)
	Update(ctx context.Context, req updateReq, userID string, eventID uuid.UUID, classID uuid.UUID) error
}

type usecase struct {
	repository *db.Queries
	db         *sql.DB
}

func NewUsecase(repository *db.Queries, db *sql.DB) Usecase {
	return &usecase{repository: repository, db: db}
}

func (u usecase) Assign(ctx context.Context, userID string, req Request, eventID uuid.UUID) error {
	eventID, err := u.repository.EventFetchOneByIDAndUserID(ctx, db.EventFetchOneByIDAndUserIDParams{
		ID:     eventID,
		UserID: uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in fetch event : %s", err.Error())
	}

	tx, err := u.db.Begin()
	if err != nil {
		return fmt.Errorf("error in start db transaction : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	for _, data := range req.Data {
		if err := txQuery.ClassEventAssign(ctx, db.ClassEventAssignParams{
			ClassID: uuid.MustParse(data.ClassID),
			EventID: eventID,
			Price:   data.Price,
		}); err != nil {
			return tx.Rollback()
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error in db commit : %s", err.Error())
	}

	return nil
}

func (u usecase) Detach(ctx context.Context, userID string, eventID uuid.UUID, classID uuid.UUID) error {
	classEventID, err := u.repository.ClassEventCheckOne(ctx, db.ClassEventCheckOneParams{
		ID:      classID,
		EventID: eventID,
		UserID:  uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("error in check class event : %s", err.Error())
	}

	res, err := u.repository.EventRegistrationFetchByClassEventIDAndStatus(ctx, classEventID)
	if err != nil {
		return fmt.Errorf("error in fetch event : %s", err.Error())
	}

	if len(res) != 0 {
		return fmt.Errorf("can't detach class event that already had approved participant")
	}

	return u.repository.ClassEventDetach(ctx, db.ClassEventDetachParams{
		ID:      classEventID,
		EventID: eventID,
	})
}

func (u usecase) FetchAll(ctx context.Context, userID string, eventID uuid.UUID) ([]db.ClassEventFetchAllRow, error) {
	return u.repository.ClassEventFetchAll(ctx, db.ClassEventFetchAllParams{
		ID:     eventID,
		UserID: uuid.MustParse(userID),
	})
}

func (u usecase) Update(ctx context.Context, req updateReq, userID string, eventID uuid.UUID, classID uuid.UUID) error {
	classEventID, err := u.repository.ClassEventCheckOne(ctx, db.ClassEventCheckOneParams{
		ID:      classID,
		EventID: eventID,
		UserID:  uuid.MustParse(userID),
	})
	if err != nil {
		return fmt.Errorf("record not found in class event : %s", err.Error())
	}
	return u.repository.ClassEventUpdate(ctx, db.ClassEventUpdateParams{
		Price: req.Price,
		ID:    classEventID,
	})
}
