package score_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/score"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_score "github.com/project-ippl-dev/tanding-api/mocks/score"
)

type handlerMock struct {
	mockUsecase *mock_score.MockUsecase
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUsecase := mock_score.NewMockUsecase(ctrl)
	e := echo.New()
	return handlerMock{
		mockUsecase: mockUsecase,
	}, e, ctrl
}

func TestHandler_StoreOrUpdateOrder(t *testing.T) {
	testCases := []struct {
		description        string
		requestBody        string
		ucErr              error
		ucStatus           int
		expectedMessage    string
		expectedStatusCode int
		setMock            func(c echo.Context, mock handlerMock)
	}{
		{
			description:        "bind error",
			requestBody:        "not-a-json",
			ucErr:              nil,
			ucStatus:           0,
			expectedMessage:    "Unmarshal",
			expectedStatusCode: http.StatusBadRequest,
			setMock:            func(c echo.Context, mock handlerMock) {},
		},
		{
			description:        "validation error",
			requestBody:        `{"order_bracket_id":"00000000-0000-0000-0000-000000000000","round_1":0,"round_2":0,"round_3":0,"extra":0,"total":-1}`,
			ucErr:              nil,
			ucStatus:           0,
			expectedMessage:    "error validation",
			expectedStatusCode: http.StatusUnprocessableEntity,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().StoreOrUpdateOrder(gomock.Any(), gomock.Any()).AnyTimes()
			},
		},
		{
			description:        "usecase error",
			requestBody:        `{"field":"valid"}`,
			ucErr:              errors.New("fail"),
			ucStatus:           http.StatusInternalServerError,
			expectedMessage:    "fail",
			expectedStatusCode: http.StatusInternalServerError,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().StoreOrUpdateOrder(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("fail"))
			},
		},
		{
			description:        "success",
			requestBody:        `{"field":"valid"}`,
			ucErr:              nil,
			ucStatus:           http.StatusCreated,
			expectedMessage:    "store or update order scores success",
			expectedStatusCode: http.StatusCreated,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().StoreOrUpdateOrder(gomock.Any(), gomock.Any()).Return(http.StatusCreated, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e, ctrl := newHandlerMock(t)
			defer ctrl.Finish()
			h := score.NewHandler(mock.mockUsecase)
			req := httptest.NewRequest(http.MethodPost, "/score/order", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)
			c.SetPath("/score/order")
			testCase.setMock(c, mock)
			err := h.StoreOrUpdateOrder(c)
			assert.NoError(t, err)
			var resp tools.Response
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandler_FetchOneOrder(t *testing.T) {
	testCases := []struct {
		description        string
		ucData             interface{}
		ucErr              error
		expectedMessage    string
		expectedStatusCode int
		setMock            func(c echo.Context, mock handlerMock)
	}{
		{
			description:        "success",
			ucData:             db.OrderScoreFetchOneByBracketIDRow{ID: uuid.New(), Round1: 1, Round2: 2, Round3: 3, Extra: 0, Total: 6},
			ucErr:              nil,
			expectedMessage:    "fetch one order score success",
			expectedStatusCode: http.StatusOK,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().FetchOneOrder(gomock.Any(), gomock.Any()).Return(db.OrderScoreFetchOneByBracketIDRow{ID: uuid.New(), Round1: 1, Round2: 2, Round3: 3, Extra: 0, Total: 6}, nil)
			},
		},
		{
			description:        "usecase error",
			ucData:             db.OrderScoreFetchOneByBracketIDRow{},
			ucErr:              errors.New("not found"),
			expectedMessage:    "not found",
			expectedStatusCode: http.StatusNotFound,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().FetchOneOrder(gomock.Any(), gomock.Any()).Return(db.OrderScoreFetchOneByBracketIDRow{}, errors.New("not found"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e, ctrl := newHandlerMock(t)
			defer ctrl.Finish()
			h := score.NewHandler(mock.mockUsecase)
			req := httptest.NewRequest(http.MethodGet, "/score/order", nil)
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)
			c.SetPath("/score/order")
			testCase.setMock(c, mock)
			err := h.FetchOneOrder(c)
			assert.NoError(t, err)
			if testCase.ucErr == nil {
				var resp tools.ResponseData
				_ = json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.Equal(t, testCase.expectedStatusCode, rr.Code)
				assert.Equal(t, testCase.expectedMessage, resp.Message)
				// assert.Equal(t, testCase.ucData, resp.Data)
			} else {
				var resp tools.Response
				_ = json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.Equal(t, testCase.expectedStatusCode, rr.Code)
				assert.Equal(t, testCase.expectedMessage, resp.Message)
			}
		})
	}
}

func TestHandler_StoreOrUpdateSingle(t *testing.T) {
	testCases := []struct {
		description        string
		requestBody        string
		ucErr              error
		ucStatus           int
		expectedMessage    string
		expectedStatusCode int
		setMock            func(c echo.Context, mock handlerMock)
	}{
		{
			description:        "bind error",
			requestBody:        "not-a-json",
			ucErr:              nil,
			ucStatus:           0,
			expectedMessage:    "Unmarshal",
			expectedStatusCode: http.StatusBadRequest,
			setMock:            func(c echo.Context, mock handlerMock) {},
		},
		{
			description:        "validation error",
			requestBody:        `{"event_bracket_id":"00000000-0000-0000-0000-000000000000","home_total":-1,"away_total":-1}`,
			ucErr:              nil,
			ucStatus:           0,
			expectedMessage:    "error validation",
			expectedStatusCode: http.StatusUnprocessableEntity,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().StoreOrUpdateSingle(gomock.Any(), gomock.Any()).AnyTimes()
			},
		},
		{
			description:        "usecase error",
			requestBody:        `{"field":"valid"}`,
			ucErr:              errors.New("fail"),
			ucStatus:           http.StatusInternalServerError,
			expectedMessage:    "fail",
			expectedStatusCode: http.StatusInternalServerError,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().StoreOrUpdateSingle(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("fail"))
			},
		},
		{
			description:        "success",
			requestBody:        `{"field":"valid"}`,
			ucErr:              nil,
			ucStatus:           http.StatusCreated,
			expectedMessage:    "store or update single scores success",
			expectedStatusCode: http.StatusCreated,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().StoreOrUpdateSingle(gomock.Any(), gomock.Any()).Return(http.StatusCreated, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e, ctrl := newHandlerMock(t)
			defer ctrl.Finish()
			h := score.NewHandler(mock.mockUsecase)
			req := httptest.NewRequest(http.MethodPost, "/score/single", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)
			c.SetPath("/score/single")
			testCase.setMock(c, mock)
			err := h.StoreOrUpdateSingle(c)
			assert.NoError(t, err)
			var resp tools.Response
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandler_FetchOneSingle(t *testing.T) {
	testCases := []struct {
		description        string
		ucData             interface{}
		ucErr              error
		expectedMessage    string
		expectedStatusCode int
		setMock            func(c echo.Context, mock handlerMock)
	}{
		{
			description:        "success",
			ucData:             db.EventScoreFetchOneByBracketIDRow{ID: uuid.New(), HomeRound1: 1, HomeRound2: 2, HomeRound3: 3, HomeExtra: 0, HomeTotal: 6, AwayRound1: 1, AwayRound2: 2, AwayRound3: 3, AwayExtra: 0, AwayTotal: 6},
			ucErr:              nil,
			expectedMessage:    "fetch one single score success",
			expectedStatusCode: http.StatusOK,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().FetchOneSingle(gomock.Any(), gomock.Any()).Return(db.EventScoreFetchOneByBracketIDRow{ID: uuid.New(), HomeRound1: 1, HomeRound2: 2, HomeRound3: 3, HomeExtra: 0, HomeTotal: 6, AwayRound1: 1, AwayRound2: 2, AwayRound3: 3, AwayExtra: 0, AwayTotal: 6}, nil)
			},
		},
		{
			description:        "usecase error",
			ucData:             db.EventScoreFetchOneByBracketIDRow{},
			ucErr:              errors.New("not found"),
			expectedMessage:    "not found",
			expectedStatusCode: http.StatusNotFound,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().FetchOneSingle(gomock.Any(), gomock.Any()).Return(db.EventScoreFetchOneByBracketIDRow{}, errors.New("not found"))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e, ctrl := newHandlerMock(t)
			defer ctrl.Finish()
			h := score.NewHandler(mock.mockUsecase)
			req := httptest.NewRequest(http.MethodGet, "/score/single", nil)
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)
			c.SetPath("/score/single")
			testCase.setMock(c, mock)
			err := h.FetchOneSingle(c)
			assert.NoError(t, err)
			if testCase.ucErr == nil {
				var resp tools.ResponseData
				_ = json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.Equal(t, testCase.expectedStatusCode, rr.Code)
				assert.Equal(t, testCase.expectedMessage, resp.Message)
				// assert.Equal(t, testCase.ucData, resp.Data)
			} else {
				var resp tools.Response
				_ = json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.Equal(t, testCase.expectedStatusCode, rr.Code)
				assert.Equal(t, testCase.expectedMessage, resp.Message)
			}
		})
	}
}

func TestHandler_Lock(t *testing.T) {
	testCases := []struct {
		description        string
		requestBody        string
		ucErr              error
		ucStatus           int
		expectedMessage    string
		expectedStatusCode int
		setMock            func(c echo.Context, mock handlerMock)
	}{
		{
			description:        "bind error",
			requestBody:        "not-a-json",
			ucErr:              nil,
			ucStatus:           0,
			expectedMessage:    "Unmarshal",
			expectedStatusCode: http.StatusBadRequest,
			setMock:            func(c echo.Context, mock handlerMock) {},
		},
		{
			description:        "usecase error",
			requestBody:        `{"field":"valid"}`,
			ucErr:              errors.New("fail"),
			ucStatus:           http.StatusInternalServerError,
			expectedMessage:    "fail",
			expectedStatusCode: http.StatusInternalServerError,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().Lock(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("fail"))
			},
		},
		{
			description:        "success",
			requestBody:        `{"field":"valid"}`,
			ucErr:              nil,
			ucStatus:           http.StatusOK,
			expectedMessage:    "lock score success",
			expectedStatusCode: http.StatusOK,
			setMock: func(c echo.Context, mock handlerMock) {
				mock.mockUsecase.EXPECT().Lock(gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e, ctrl := newHandlerMock(t)
			defer ctrl.Finish()
			h := score.NewHandler(mock.mockUsecase)
			req := httptest.NewRequest(http.MethodPost, "/score/lock", strings.NewReader(testCase.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			c := e.NewContext(req, rr)
			c.SetPath("/score/lock")
			testCase.setMock(c, mock)
			err := h.Lock(c)
			assert.NoError(t, err)
			var resp tools.Response
			_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		})
	}
}
