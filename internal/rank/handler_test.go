package rank_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/google/uuid"
	"github.com/project-ippl-dev/tanding-api/internal/rank"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_rank "github.com/project-ippl-dev/tanding-api/mocks/rank"
	"github.com/project-ippl-dev/tanding-api/testutils"
)

type mockJWTClient struct {
	ID string
}

func (m *mockJWTClient) Decode(c echo.Context) struct{ ID string } {
	return struct{ ID string }{ID: m.ID}
}

type handlerMock struct {
	mockUsecase *mock_rank.MockUsecase
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)

	mockRankUsecase := mock_rank.NewMockUsecase(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase: mockRankUsecase,
	}, e
}
func TestHandler_summary(t *testing.T) {
	mock, e := newHandlerMock(t)
	fmt.Println(mock)

	h := rank.NewHandler(mock.mockUsecase, tools.JWTClient)

	eventID := uuid.New()

	httpMethod := http.MethodPost
	url := "/event/:event/summary"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedMessage    string
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "bind error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "code=400, message=invalid UUID length: 22, internal=invalid UUID length: 22",
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("event")
				c.SetParamValues("invalid-event-id")
			},
		},
		{
			description: "usecase summary return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "error",
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("event")
				c.SetParamValues(eventID.String())
				mock.mockUsecase.EXPECT().Summary(gomock.Any(), eventID).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "generate summary event success",
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("event")
				c.SetParamValues(eventID.String())
				mock.mockUsecase.EXPECT().Summary(gomock.Any(), eventID).Return(http.StatusCreated, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := h.summary(c)
		assert.NoError(t, err)

		var response tools.Response
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		assert.Equal(t, testCase.expectedMessage, response.Message)
	}
}

func TestHandler_FetchOwnPoint(t *testing.T) {
	mock, e := newHandlerMock(t)

	userID := uuid.New().String()
	h := rank.NewHandler(mock.mockUsecase, tools.JWTClient)

	httpMethod := http.MethodGet
	url := "/rank/point/own"

	testCases := []struct {
		description        string
		point              int32
		error              error
		expectedMessage    string
		expectedStatusCode int
	}{
		{
			description:        "success",
			point:              10,
			error:              nil,
			expectedMessage:    "fetch point by auth user success",
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "usecase error",
			point:              0,
			error:              fmt.Errorf("not found"),
			expectedMessage:    "not found",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testutils.MockHttpRequestParam{
			HttpMethod: httpMethod,
			Url:        url,
		})
		c := e.NewContext(httpReq, rr)
		mock.mockUsecase.EXPECT().FetchOwnPoint(gomock.Any(), userID).Return(testCase.point, testCase.error)
		err := h.FetchOwnPoint(c)
		assert.NoError(t, err)
		if testCase.error == nil {
			var response tools.ResponseData
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
			assert.Equal(t, testCase.point, response.Data)
		} else {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
		}
	}
}

func TestHandler_FetchByClubID(t *testing.T) {
	mock, e := newHandlerMock(t)

	h := rank.NewHandler(mock.mockUsecase, tools.JWTClient)

	clubID := uuid.New()
	url := "/rank/point/club/:club"
	httpMethod := http.MethodGet

	fakeResp := map[string]interface{}{"TotalPoint": 10, "Participants": []interface{}{"user1", "user2"}}

	testCases := []struct {
		description        string
		clubID             string
		mockResp           interface{}
		mockErr            error
		expectedMessage    string
		expectedStatusCode int
	}{
		{
			description:        "bind error",
			clubID:             "invalid-uuid",
			mockResp:           nil,
			mockErr:            nil,
			expectedMessage:    "invalid UUID length: 12",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "usecase error",
			clubID:             clubID.String(),
			mockResp:           nil,
			mockErr:            fmt.Errorf("not found"),
			expectedMessage:    "not found",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description:        "success",
			clubID:             clubID.String(),
			mockResp:           fakeResp,
			mockErr:            nil,
			expectedMessage:    "fetch point by club id success",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testutils.MockHttpRequestParam{
			HttpMethod: httpMethod,
			Url:        url,
		})
		c := e.NewContext(httpReq, rr)
		c.SetPath(url)
		c.SetParamNames("club")
		c.SetParamValues(testCase.clubID)
		if testCase.description == "success" || testCase.description == "usecase error" {
			parsedID, _ := uuid.Parse(testCase.clubID)
			mock.mockUsecase.EXPECT().FetchByClubID(gomock.Any(), parsedID).Return(testCase.mockResp, testCase.mockErr)
		}
		err := h.FetchByClubID(c)
		assert.NoError(t, err)
		if testCase.description == "success" {
			var response tools.ResponseData
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
			assert.Equal(t, testCase.mockResp, response.Data)
		} else {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
		}
	}
}

func TestHandler_Rank(t *testing.T) {
	mock, e := newHandlerMock(t)

	h := rank.NewHandler(mock.mockUsecase, tools.JWTClient)

	httpMethod := http.MethodGet
	url := "/rank/club"

	fakePagination := tools.Pagination{TotalItem: 1, PageSize: 10, Page: 1, Data: []interface{}{"user1"}}

	testCases := []struct {
		description        string
		bindErr            bool
		validateErr        bool
		mockResp           tools.Pagination
		mockErr            error
		expectedMessage    string
		expectedStatusCode int
	}{
		{
			description:        "bind error",
			bindErr:            true,
			validateErr:        false,
			mockResp:           tools.Pagination{},
			mockErr:            nil,
			expectedMessage:    "code=400, message=bind error",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "validation error",
			bindErr:            false,
			validateErr:        true,
			mockResp:           tools.Pagination{},
			mockErr:            nil,
			expectedMessage:    "error validation",
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			description:        "usecase error",
			bindErr:            false,
			validateErr:        false,
			mockResp:           tools.Pagination{},
			mockErr:            fmt.Errorf("internal error"),
			expectedMessage:    "internal error",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description:        "success",
			bindErr:            false,
			validateErr:        false,
			mockResp:           fakePagination,
			mockErr:            nil,
			expectedMessage:    "fetch all rank success",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testutils.MockHttpRequestParam{
			HttpMethod: httpMethod,
			Url:        url,
		})
		c := e.NewContext(httpReq, rr)
		if testCase.bindErr {
			c.Set("bindError", true)
		} else if testCase.validateErr {
			c.Set("validateError", true)
		} else {
			mock.mockUsecase.EXPECT().Rank(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.mockResp, testCase.mockErr)
		}
		err := h.rank(c)
		assert.NoError(t, err)
		if testCase.description == "success" {
			var response tools.ResponseData
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
			assert.Equal(t, testCase.mockResp, response.Data)
		} else if testCase.description == "validation error" {
			var response tools.ResponseValidation
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
		} else {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
		}
	}
}

func TestHandler_UserRank(t *testing.T) {
	mock, e := newHandlerMock(t)

	h := rank.NewHandler(mock.mockUsecase, tools.JWTClient)

	httpMethod := http.MethodGet
	url := "/rank/user"

	fakePagination := tools.Pagination{TotalItem: 1, PageSize: 10, Page: 1, Data: []interface{}{"user1"}}

	testCases := []struct {
		description        string
		bindErr            bool
		validateErr        bool
		mockResp           tools.Pagination
		mockErr            error
		expectedMessage    string
		expectedStatusCode int
	}{
		{
			description:        "bind error",
			bindErr:            true,
			validateErr:        false,
			mockResp:           tools.Pagination{},
			mockErr:            nil,
			expectedMessage:    "code=400, message=bind error",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "validation error",
			bindErr:            false,
			validateErr:        true,
			mockResp:           tools.Pagination{},
			mockErr:            nil,
			expectedMessage:    "error validation",
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			description:        "usecase error",
			bindErr:            false,
			validateErr:        false,
			mockResp:           tools.Pagination{},
			mockErr:            fmt.Errorf("internal error"),
			expectedMessage:    "internal error",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description:        "success",
			bindErr:            false,
			validateErr:        false,
			mockResp:           fakePagination,
			mockErr:            nil,
			expectedMessage:    "fetch all user rank success",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testutils.MockHttpRequestParam{
			HttpMethod: httpMethod,
			Url:        url,
		})
		c := e.NewContext(httpReq, rr)
		if testCase.bindErr {
			c.Set("bindError", true)
		} else if testCase.validateErr {
			c.Set("validateError", true)
		} else {
			mock.mockUsecase.EXPECT().UserRank(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.mockResp, testCase.mockErr)
		}
		err := h.userRank(c)
		assert.NoError(t, err)
		if testCase.description == "success" {
			var response tools.ResponseData
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
			assert.Equal(t, testCase.mockResp, response.Data)
		} else if testCase.description == "validation error" {
			var response tools.ResponseValidation
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
		} else {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
		}
	}
}
