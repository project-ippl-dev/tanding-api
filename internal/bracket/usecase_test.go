package bracket_test

import (
	"context"
	"fmt"
	"github.com/project-ippl-dev/tanding-api/internal/bracket"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	mock_bracket "github.com/project-ippl-dev/tanding-api/mocks/bracket"
	bracketFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/bracket"
	dbFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/db"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

type usecaseMock struct {
	mockRawRepository *mock_bracket.MockRepository
	mockSqlDB         sqlmock.Sqlmock
}

func newUsecaseMock(t *testing.T) (context.Context, usecaseMock, bracket.Usecase) {
	ctx := context.Background()
	mockCtrl, mockCtx := gomock.WithContext(ctx, t)
	defer mockCtrl.Finish()

	mockRawRepository := mock_bracket.NewMockRepository(mockCtrl)

	sqlDB, mockSqlDB, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repository := db.New(sqlDB)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	bracketUsecase := bracket.NewUsecase(repository, mockRawRepository, r, sqlDB)

	return mockCtx, usecaseMock{
		mockRawRepository: mockRawRepository,
		mockSqlDB:         mockSqlDB,
	}, bracketUsecase
}

func TestUsecase_Store(t *testing.T) {
	mockCtx, mock, bracketUsecase := newUsecaseMock(t)

	req := bracket.GenerateParams{
		ClassEventID: bracketFixtures.ClassEventID,
	}

	mockReqClassEventGeneratedTypeOrder := dbFixtures.MockResponseClassEventFetchOneRowReq{
		MatchType:         db.MatchTypeOrder,
		IsBracketGenerate: true,
		RuleMale:          1,
		RuleFemale:        0,
		RuleTotal:         1,
		MatchIndex:        1,
	}

	validMockReqClassEventTypeOrder := dbFixtures.MockResponseClassEventFetchOneRowReq{
		MatchType:  db.MatchTypeOrder,
		RuleMale:   1,
		RuleFemale: 0,
		RuleTotal:  1,
		MatchIndex: 1,
	}

	validMockReqClassEventTypeSingle := dbFixtures.MockResponseClassEventFetchOneRowReq{
		MatchType:  db.MatchTypeSingle,
		RuleMale:   1,
		RuleFemale: 0,
		RuleTotal:  1,
		MatchIndex: 1,
	}

	validMockReqEventRegistrationFetchByClassEventIDTypeOrder := dbFixtures.MockResponseEventRegistrationFetchByClassEventIDReq{
		RegistrationIteration: 4,
	}

	validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne := dbFixtures.MockResponseEventRegistrationFetchByClassEventIDReq{
		RegistrationIteration: 1,
	}

	//validMockReqEventRegistrationFetchByClassEventIDTypeSingleFour := dbFixtures.MockResponseEventRegistrationFetchByClassEventIDReq{
	//	RegistrationIteration: 4,
	//}

	testCases := []struct {
		description string
		req         bracket.GenerateParams
		expectedErr error
		expectedRes int
		beforeTest  func(mock usecaseMock)
	}{
		{
			description: "repository ClassEventFetchOne return error",
			req:         req,
			expectedErr: fmt.Errorf("class event not found : %s", "error"),
			expectedRes: http.StatusNotFound,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "classEvent already generated before",
			req:         req,
			expectedErr: fmt.Errorf("bracket for specific class event already generated"),
			expectedRes: http.StatusForbidden,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(mockReqClassEventGeneratedTypeOrder))
			},
		},
		{
			description: "repository EventRegistrationFetchByClassEventID return error",
			req:         req,
			expectedErr: fmt.Errorf("registration not found : %s", "error"),
			expectedRes: http.StatusNotFound,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "running db tx return error",
			req:         req,
			expectedErr: fmt.Errorf("error in start tx : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeOrder))
				mock.mockSqlDB.
					ExpectBegin().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "repository ClassEventUpdateBracketGenerate return error and fail rollback",
			req:         req,
			expectedErr: fmt.Errorf("error in rollback tx update class event : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeOrder))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "repository ClassEventUpdateBracketGenerate return error",
			req:         req,
			expectedErr: fmt.Errorf("error in update class event : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeOrder))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback()
			},
		},
		{
			description: "class event type order, repository OrderBracketCreate return error and fail rollback",
			req:         req,
			expectedErr: fmt.Errorf("error in rollback tx generate bracket order : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeOrder))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "class event type order, repository OrderBracketCreate return error",
			req:         req,
			expectedErr: fmt.Errorf("error in generate bracket order : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeOrder))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback()
			},
		},
		{
			description: "class event type order, db commit return error",
			req:         req,
			expectedErr: fmt.Errorf("error in commit tx : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeOrder))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration)))
				for i := range validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration {
					mock.mockSqlDB.
						ExpectExec(".*").
						WillReturnResult(sqlmock.NewResult(int64(i), int64(i)))
				}
				mock.mockSqlDB.
					ExpectCommit().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "class event type order, success",
			req:         req,
			expectedErr: nil,
			expectedRes: http.StatusCreated,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeOrder))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeOrder))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration)))
				for i := range validMockReqEventRegistrationFetchByClassEventIDTypeOrder.RegistrationIteration {
					mock.mockSqlDB.
						ExpectExec(".*").
						WillReturnResult(sqlmock.NewResult(int64(i), int64(i)))
				}
				mock.mockSqlDB.
					ExpectCommit()
			},
		},
		{
			description: "class event type single with one participant, repository EventBracketCreate return error and fail rollback",
			req:         req,
			expectedErr: fmt.Errorf("error in rollback tx create event bracket : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "class event type single with one participant, repository EventBracketCreate return error",
			req:         req,
			expectedErr: fmt.Errorf("error in create event bracket : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback()
			},
		},
		{
			description: "class event type single with one participant, repository BracketParticipantCreate for Home return error and fail rollback",
			req:         req,
			expectedErr: fmt.Errorf("error in rollback tx create bracket participant : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventBracketCreate(bracketFixtures.BracketID))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "class event type single with one participant, repository BracketParticipantCreate for Home return error",
			req:         req,
			expectedErr: fmt.Errorf("error in create bracket participant : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventBracketCreate(bracketFixtures.BracketID))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback()
			},
		},
		{
			description: "class event type single with one participant, repository BracketParticipantCreate for Away return error and fail rollback",
			req:         req,
			expectedErr: fmt.Errorf("error in rollback tx create bracket participant : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventBracketCreate(bracketFixtures.BracketID))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "class event type single with one participant, repository BracketParticipantCreate for Away return error",
			req:         req,
			expectedErr: fmt.Errorf("error in create bracket participant : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventBracketCreate(bracketFixtures.BracketID))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback()
			},
		},
		{
			description: "class event type single with one participant, success",
			req:         req,
			expectedErr: nil,
			expectedRes: http.StatusCreated,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventBracketCreate(bracketFixtures.BracketID))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.mockSqlDB.
					ExpectCommit()
			},
		},
		{
			description: "class event type single with four participant, repository EventBracketCreate return error and fail rollback",
			req:         req,
			expectedErr: fmt.Errorf("error in rollback tx create event bracket : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseClassEventFetchOneRow(validMockReqClassEventTypeSingle))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(dbFixtures.NewMockResponseEventRegistrationFetchByClassEventID(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne))
				mock.mockSqlDB.
					ExpectBegin()
				mock.mockSqlDB.
					ExpectExec(".*").
					WillReturnResult(sqlmock.NewResult(int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration), int64(validMockReqEventRegistrationFetchByClassEventIDTypeSingleOne.RegistrationIteration)))
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnError(fmt.Errorf("error"))
				mock.mockSqlDB.
					ExpectRollback().
					WillReturnError(fmt.Errorf("error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.beforeTest(mock)
			res, err := bracketUsecase.Store(mockCtx, tc.req)
			assert.Equal(t, tc.expectedRes, res, "res must be %t", tc.expectedRes)
			assert.Equal(t, tc.expectedErr, err, "err must be %t", tc.expectedErr)
		})
	}
}
