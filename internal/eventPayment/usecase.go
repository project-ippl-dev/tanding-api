package eventPayment

//go:generate mockgen -source=./usecase.go -destination=../../mocks/eventPayment/usecase_mock.go

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Usecase interface {
	Store(ctx context.Context, req Request, userID string) (statusCode int, err error)
	FetchAll(ctx context.Context, page int32, pageSize int32, arg FetchAllParams) (tools.Pagination, error)
	Update(ctx context.Context, arg UpdateParams, userID string) (statusCode int, err error)
	FetchByUserID(ctx context.Context, page, pageSize int32, arg FetchByUserIDParams, userID string) (statusCode int, pagination tools.Pagination, err error)
	Cart(ctx context.Context, page, pageSize int32, userID string) (tools.Pagination, error)
	CartDetail(ctx context.Context, eventID uuid.UUID, userID string) (CartDetailResponse, error)
	Detail(ctx context.Context, paymentID uuid.UUID) (DetailPaymentResponse, error)
	Summary(ctx context.Context, arg SummaryParams) (SummaryResponse, error)
}

type usecase struct {
	repository    *db.Queries
	rawRepository RawRepository
	rdb           *redis.Client
	r             *rand.Rand
}

func NewUsecase(repository *db.Queries, rawRepository RawRepository, rdb *redis.Client, r *rand.Rand) Usecase {
	return &usecase{
		repository:    repository,
		rawRepository: rawRepository,
		rdb:           rdb,
		r:             r,
	}
}

func (u usecase) Store(ctx context.Context, req Request, userID string) (statusCode int, err error) {
	event, err := u.repository.EventFetchOne(ctx, req.EventID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("event not found : %s", err.Error())
	}
	if event.Deadline.Before(time.Now()) {
		return http.StatusForbidden, fmt.Errorf("register already passed. You can't continue the payment")
	}
	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	paymentID, err := txQuery.EventPaymentCreate(ctx, db.EventPaymentCreateParams{
		UniqueNumber: req.UniqueNumber,
		EventID:      req.EventID,
		UserID:       uuid.MustParse(userID),
		PaymentLink:  req.Link,
		AdminID:      uuid.UUID{},
		Total:        req.Total,
	})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx create payment : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in store payment : %s", err.Error())
	}

	var total int32 = 0

	for _, registrationID := range req.Registrations {
		registration, err := txQuery.EventRegistrationFetchOne(ctx, registrationID)
		if err != nil {
			return http.StatusNotFound, fmt.Errorf("registration record not found : %s", err.Error())
		}

		if registration.Status != db.EventRegistrationStatusPending {
			return http.StatusBadRequest, fmt.Errorf("only fetch only registration with status pending")
		}
		if err := txQuery.EventRegistrationUpdateOnePaymentID(ctx, db.EventRegistrationUpdateOnePaymentIDParams{
			EventPaymentReceiptID: paymentID,
			Status:                db.EventRegistrationStatusWaiting,
			ID:                    registrationID,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				return http.StatusInternalServerError, fmt.Errorf("error in rollback tx store payment id in specific registration : %s", err.Error())
			}
			return http.StatusInternalServerError, fmt.Errorf("error in update registration : %s", err.Error())
		}

		classEvent, err := u.repository.ClassEventFetchOne(ctx, registration.ClassEventID)
		if err != nil {
			return http.StatusNotFound, fmt.Errorf("error in fetch class event : %s", err.Error())
		}

		total += classEvent.Price
	}

	if total != req.Total {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx when total is different calculation%s", err.Error())
		}
		return http.StatusBadRequest, fmt.Errorf("total balance is difference between input and calculation")
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return http.StatusOK, nil
}

func (u usecase) FetchAll(ctx context.Context, page int32, pageSize int32, arg FetchAllParams) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)

	start, _ := time.Parse(time.RFC3339, arg.Start)
	end, _ := time.Parse(time.RFC3339, arg.End)

	queryParam := queryPaymentParams{EventID: arg.EventID,
		Status: arg.Status,
		ClubID: arg.Clubs,
		Start:  start,
		End:    end,
	}

	payments, err := u.rawRepository.EventPaymentFetchAll(ctx, fetchByEventIDParams{
		queryPaymentParams: queryParam,
		Limit:              pageSize,
		Offset:             skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch payment : %s", err.Error())
	}

	count, err := u.rawRepository.EventPaymentCountByEventID(ctx, queryParam)
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count payment : %s", err.Error())
	}

	results := []response{}
	for _, payment := range payments {
		classEvents, err := u.repository.ClassEventFetchByPaymentID(ctx, payment.ID)
		if err != nil {
			return tools.Pagination{}, fmt.Errorf("error in fetch class events : %s", err.Error())
		}
		results = append(results, response{
			fetchAllRow: payment,
			ClassEvents: classEvents,
		})
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      results,
	}, nil

}

func (u usecase) Update(ctx context.Context, arg UpdateParams, userID string) (statusCode int, err error) {
	paymentID, err := u.repository.EventPaymentCheckOne(ctx, db.EventPaymentCheckOneParams{
		EventID: arg.EventID,
		ID:      arg.PaymentID,
	})
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("record not found in payment: %s", err.Error())
	}

	tx, err := u.rawRepository.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in start tx : %s", err.Error())
	}
	txQuery := u.repository.WithTx(tx)
	if err := txQuery.EventPaymentUpdate(ctx, db.EventPaymentUpdateParams{
		Status:  arg.Status,
		ID:      paymentID,
		AdminID: uuid.MustParse(userID),
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update payment : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update payment receipt : %s", err.Error())
	}

	var status db.EventRegistrationStatus
	switch arg.Status {
	case db.EventReceiptStatusApproved:
		status = db.EventRegistrationStatusApproved
	case db.EventReceiptStatusRejected:
		status = db.EventRegistrationStatusRejected
	case db.EventReceiptStatusRefund:
		status = db.EventRegistrationStatusCanceled
	default:
		return http.StatusBadRequest, fmt.Errorf("bad request")
	}

	if err := txQuery.EventRegistrationUpdatePaymentID(ctx, db.EventRegistrationUpdatePaymentIDParams{
		EventPaymentReceiptID: paymentID,
		Status:                status,
	}); err != nil {
		if err := tx.Rollback(); err != nil {
			return http.StatusInternalServerError, fmt.Errorf("error in rollback tx update status registration : %s", err.Error())
		}
		return http.StatusInternalServerError, fmt.Errorf("error in update status registration : %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error in commit tx : %s", err.Error())
	}

	return http.StatusOK, nil
}

func (u usecase) FetchByUserID(ctx context.Context, page, pageSize int32, arg FetchByUserIDParams, userID string) (statusCode int, pagination tools.Pagination, err error) {
	skip := tools.PaginationSkip(page, pageSize)

	start, _ := time.Parse(time.RFC3339, arg.Start)
	end, _ := time.Parse(time.RFC3339, arg.End)

	clubs, err := u.repository.ClubFetchOwnerByEventID(ctx, db.ClubFetchOwnerByEventIDParams{
		EventID: arg.EventID,
		UserID:  uuid.MustParse(userID),
	})
	if err != nil {
		return http.StatusBadRequest, tools.Pagination{}, fmt.Errorf("only club owner participate in event can access : %s", err.Error())
	}

	var clubsStr []string
	for _, club := range clubs {
		clubsStr = append(clubsStr, club.String())
	}

	queryParam := queryPaymentParams{EventID: arg.EventID.String(),
		Status: arg.Status,
		ClubID: clubsStr,
		Start:  start,
		End:    end,
	}

	payments, err := u.rawRepository.EventPaymentFetchAll(ctx, fetchByEventIDParams{
		queryPaymentParams: queryParam,
		Limit:              pageSize,
		Offset:             skip,
	})

	if err != nil {
		return http.StatusInternalServerError, tools.Pagination{}, fmt.Errorf("error in fetch payment : %s", err.Error())
	}

	count, err := u.rawRepository.EventPaymentCountByEventID(ctx, queryParam)
	if err != nil {
		return http.StatusInternalServerError, tools.Pagination{}, fmt.Errorf("error in count payment : %s", err.Error())
	}

	results := []response{}
	for _, payment := range payments {
		classEvents, err := u.repository.ClassEventFetchByPaymentID(ctx, payment.ID)
		if err != nil {
			return http.StatusNotFound, tools.Pagination{}, fmt.Errorf("error in fetch class events : %s", err.Error())
		}
		results = append(results, response{
			fetchAllRow: payment,
			ClassEvents: classEvents,
		})
	}

	return http.StatusOK, tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      results,
	}, nil
}

func (u usecase) Cart(ctx context.Context, page, pageSize int32, userID string) (tools.Pagination, error) {
	skip := tools.PaginationSkip(page, pageSize)

	carts, err := u.repository.EventRegistrationFetchCart(ctx, db.EventRegistrationFetchCartParams{
		UserID: uuid.MustParse(userID),
		Limit:  pageSize,
		Offset: skip,
	})
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in fetch cart : %s", err.Error())
	}

	count, err := u.repository.EventRegistrationCountCart(ctx, uuid.MustParse(userID))
	if err != nil {
		return tools.Pagination{}, fmt.Errorf("error in count cart : %s", err.Error())
	}

	return tools.Pagination{
		TotalItem: count,
		PageSize:  pageSize,
		Page:      page,
		Data:      carts,
	}, nil
}

func (u usecase) CartDetail(ctx context.Context, eventID uuid.UUID, userID string) (CartDetailResponse, error) {
	results, err := u.rawRepository.EventRegistrationFetchCartDetails(ctx, EventRegistrationFetchCartDetailsParams{
		UserID:  uuid.MustParse(userID),
		EventID: eventID,
	})
	if err != nil {
		return CartDetailResponse{}, fmt.Errorf("error in fetch cart details : %s", err.Error())
	}
	event, err := u.repository.EventFetchForCart(ctx, eventID)
	if err != nil {
		return CartDetailResponse{}, fmt.Errorf("error in fetch event : %s", err.Error())
	}
	uniqueNumber, err := u.rdb.Get(ctx, "unique-"+userID).Result()
	if err != nil {
		u.r.Seed(time.Now().UnixNano())
		min, max := 100, 999
		uniqueNumber = fmt.Sprintf("%d", u.r.Intn(max-min+1)+min)
		if err := u.rdb.Set(ctx, "unique-"+userID, uniqueNumber, 8*time.Hour).Err(); err != nil {
			return CartDetailResponse{}, fmt.Errorf("error in set rdb unique key : %s", err.Error())
		}
	}
	timestamp, err := u.rdb.Get(ctx, "timestamp-"+userID).Result()
	if err != nil {
		now := time.Now().Format("2006-01-02T15:04:05-07")
		if err := u.rdb.Set(ctx, "timestamp-"+userID, now, 8*time.Hour).Err(); err != nil {
			return CartDetailResponse{}, fmt.Errorf("error in set rdb timestamp key : %s", err.Error())
		}
		timestamp = now
	}

	return CartDetailResponse{
		Results: results,
		Event:   event,
		UniqueNumber: UniqueNumberData{
			Number: uniqueNumber,
			Time:   timestamp,
		},
	}, nil
}

func (u usecase) Detail(ctx context.Context, paymentID uuid.UUID) (DetailPaymentResponse, error) {
	payment, err := u.repository.EventPaymentFetchOneForAdmin(ctx, paymentID)
	if err != nil {
		return DetailPaymentResponse{}, fmt.Errorf("payment not found : %s", err.Error())
	}
	details, err := u.rawRepository.EventRegistrationFetchCartDetails(ctx, EventRegistrationFetchCartDetailsParams{
		UserID:  uuid.UUID{},
		EventID: payment.EventID,
	})
	if err != nil {
		return DetailPaymentResponse{}, fmt.Errorf("detail not found : %s", err.Error())
	}
	return DetailPaymentResponse{
		Detail:      payment,
		ClassEvents: details,
	}, nil
}

func (u usecase) Summary(ctx context.Context, arg SummaryParams) (SummaryResponse, error) {
	approved, _ := u.rawRepository.EventPaymentSumAll(ctx, sumAllParams{
		EventID: arg.EventID,
		Status:  db.EventReceiptStatusApproved,
	})
	waiting, _ := u.rawRepository.EventPaymentSumAll(ctx, sumAllParams{
		EventID: arg.EventID,
		Status:  db.EventReceiptStatusWaiting,
	})
	refund, _ := u.rawRepository.EventPaymentSumAll(ctx, sumAllParams{
		EventID: arg.EventID,
		Status:  db.EventReceiptStatusRefund,
	})
	return SummaryResponse{
		TotalApproved: approved,
		TotalWaiting:  waiting,
		TotalRefund:   refund,
	}, nil
}
