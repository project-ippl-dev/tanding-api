package rank_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/rank"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_rank "github.com/project-ippl-dev/tanding-api/mocks/rank"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type usecaseMock struct {
	mockRawRepository *mock_rank.MockRepository
	mockSqlDB         sqlmock.Sqlmock
}

func newUsecaseMock(t *testing.T) (context.Context, usecaseMock, rank.Usecase) {
	ctx := context.Background()
	mockCtrl, mockCtx := gomock.WithContext(ctx, t)
	defer mockCtrl.Finish()

	mockRawRepository := mock_rank.NewMockRepository(mockCtrl)

	sqlDB, mockSqlDB, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	repository := db.New(sqlDB)

	rankUsecase := rank.NewUsecase(repository, mockRawRepository)

	return mockCtx, usecaseMock{
		mockRawRepository: mockRawRepository,
		mockSqlDB:         mockSqlDB,
	}, rankUsecase
}

func TestUsecase_Summary(t *testing.T) {
	mockCtx, mock, rankUsecase := newUsecaseMock(t)

	eventID := uuid.New()

	testCases := []struct {
		description string
		eventID     uuid.UUID
		expectedErr error
		expectedRes int
		beforeTest  func(mock usecaseMock)
	}{
		{
			description: "running db tx return error",
			eventID:     eventID,
			expectedErr: fmt.Errorf("error in start tx : %s", "error"),
			expectedRes: http.StatusInternalServerError,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectBegin().
					WillReturnError(fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			eventID:     eventID,
			expectedErr: nil,
			expectedRes: http.StatusCreated,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.ExpectBegin()
				mock.mockSqlDB.ExpectCommit()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.beforeTest(mock)
			res, err := rankUsecase.Summary(mockCtx, tc.eventID)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestUsecase_FetchOwnPoint(t *testing.T) {
	mockCtx, mock, rankUsecase := newUsecaseMock(t)

	userID := uuid.New().String()

	testCases := []struct {
		description  string
		userID       string
		expectedErr  error
		expectedResp int32
		beforeTest   func(mock usecaseMock)
	}{
		{
			description:  "repository error",
			userID:       userID,
			expectedErr:  fmt.Errorf("failed to fetch point"),
			expectedResp: 0,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnError(fmt.Errorf("failed to fetch point"))
			},
		},
		{
			description:  "success",
			userID:       userID,
			expectedErr:  nil,
			expectedResp: 100,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(sqlmock.NewRows([]string{"point"}).AddRow(100))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.beforeTest(mock)
			point, err := rankUsecase.FetchOwnPoint(mockCtx, tc.userID)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResp, point)
		})
	}
}

func TestUsecase_FetchByClubID(t *testing.T) {
	mockCtx, mock, rankUsecase := newUsecaseMock(t)

	clubID := uuid.New()
	resp := map[string]interface{}{
		"TotalPoint":   100,
		"Participants": []string{"user1", "user2"},
	}

	testCases := []struct {
		description  string
		clubID       uuid.UUID
		expectedErr  error
		expectedResp interface{}
		beforeTest   func(mock usecaseMock)
	}{
		{
			description:  "repository error",
			clubID:       clubID,
			expectedErr:  fmt.Errorf("failed to fetch club data"),
			expectedResp: nil,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnError(fmt.Errorf("failed to fetch club data"))
			},
		},
		{
			description:  "success",
			clubID:       clubID,
			expectedErr:  nil,
			expectedResp: resp,
			beforeTest: func(mock usecaseMock) {
				mock.mockSqlDB.
					ExpectQuery(".*").
					WillReturnRows(sqlmock.NewRows([]string{"data"}).AddRow(resp))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.beforeTest(mock)
			result, err := rankUsecase.FetchByClubID(mockCtx, tc.clubID)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedResp, result)
		})
	}
}

func TestUsecase_Rank(t *testing.T) {
	mockCtx, mock, rankUsecase := newUsecaseMock(t)

	fakePagination := tools.Pagination{
		TotalItem: 1,
		PageSize:  10,
		Page:      1,
		Data:      []interface{}{"rank1"},
	}

	testCases := []struct {
		description string
		page        int32
		pageSize    int32
		arg         interface{}
		expectedErr error
		expectedRes tools.Pagination
		beforeTest  func(mock usecaseMock)
	}{
		{
			description: "repository RankFetchPowerList error",
			page:        1,
			pageSize:    10,
			arg:         "sport1",
			expectedErr: fmt.Errorf("failed to fetch power list"),
			expectedRes: tools.Pagination{},
			beforeTest: func(mock usecaseMock) {
				mock.mockRawRepository.EXPECT().
					RankFetchPowerList(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("failed to fetch power list"))
			},
		},
		{
			description: "repository RankCountPowerList error",
			page:        1,
			pageSize:    10,
			arg:         "sport1",
			expectedErr: fmt.Errorf("failed to count power list"),
			expectedRes: tools.Pagination{},
			beforeTest: func(mock usecaseMock) {
				mock.mockRawRepository.EXPECT().
					RankFetchPowerList(gomock.Any(), gomock.Any()).
					Return([]interface{}{"rank1"}, nil)
				mock.mockRawRepository.EXPECT().
					RankCountPowerList(gomock.Any(), gomock.Any()).
					Return(int64(0), fmt.Errorf("failed to count power list"))
			},
		},
		{
			description: "success",
			page:        1,
			pageSize:    10,
			arg:         "sport1",
			expectedErr: nil,
			expectedRes: fakePagination,
			beforeTest: func(mock usecaseMock) {
				mock.mockRawRepository.EXPECT().
					RankFetchPowerList(gomock.Any(), gomock.Any()).
					Return([]interface{}{"rank1"}, nil)
				mock.mockRawRepository.EXPECT().
					RankCountPowerList(gomock.Any(), gomock.Any()).
					Return(int64(1), nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.beforeTest(mock)
			result, err := rankUsecase.Rank(mockCtx, tc.page, tc.pageSize, tc.arg)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedRes, result)
		})
	}
}

func TestUsecase_UserRank(t *testing.T) {
	mockCtx, mock, rankUsecase := newUsecaseMock(t)

	fakePagination := tools.Pagination{
		TotalItem: 1,
		PageSize:  10,
		Page:      1,
		Data:      []interface{}{"user1"},
	}

	testCases := []struct {
		description string
		page        int32
		pageSize    int32
		arg         interface{}
		expectedErr error
		expectedRes tools.Pagination
		beforeTest  func(mock usecaseMock)
	}{
		{
			description: "repository RankFetchAllPointUser error",
			page:        1,
			pageSize:    10,
			arg:         "sport1",
			expectedErr: fmt.Errorf("failed to fetch user points"),
			expectedRes: tools.Pagination{},
			beforeTest: func(mock usecaseMock) {
				mock.mockRawRepository.EXPECT().
					RankFetchAllPointUser(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("failed to fetch user points"))
			},
		},
		{
			description: "repository RankCountAllPointUser error",
			page:        1,
			pageSize:    10,
			arg:         "sport1",
			expectedErr: fmt.Errorf("failed to count user points"),
			expectedRes: tools.Pagination{},
			beforeTest: func(mock usecaseMock) {
				mock.mockRawRepository.EXPECT().
					RankFetchAllPointUser(gomock.Any(), gomock.Any()).
					Return([]interface{}{"user1"}, nil)
				mock.mockRawRepository.EXPECT().
					RankCountAllPointUser(gomock.Any(), gomock.Any()).
					Return(int64(0), fmt.Errorf("failed to count user points"))
			},
		},
		{
			description: "success",
			page:        1,
			pageSize:    10,
			arg:         "sport1",
			expectedErr: nil,
			expectedRes: fakePagination,
			beforeTest: func(mock usecaseMock) {
				mock.mockRawRepository.EXPECT().
					RankFetchAllPointUser(gomock.Any(), gomock.Any()).
					Return([]interface{}{"user1"}, nil)
				mock.mockRawRepository.EXPECT().
					RankCountAllPointUser(gomock.Any(), gomock.Any()).
					Return(int64(1), nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tc.beforeTest(mock)
			result, err := rankUsecase.UserRank(mockCtx, tc.page, tc.pageSize, tc.arg)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedRes, result)
		})
	}
}
