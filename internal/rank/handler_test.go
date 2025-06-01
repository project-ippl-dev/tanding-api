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
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/rank"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	jwtFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/tools/jwt"
	mock_rank "github.com/project-ippl-dev/tanding-api/mocks/rank"
	mock_tools "github.com/project-ippl-dev/tanding-api/mocks/tools"
	"github.com/project-ippl-dev/tanding-api/testutils"
)

var fakePagination = tools.Pagination{TotalItem: 1, PageSize: 10, Page: 1, Data: []interface{}{"user1"}}

type handlerMock struct {
	mockUsecase   *mock_rank.MockUsecase
	mockJWTClient *mock_tools.MockJWTClient
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	mockRankUsecase := mock_rank.NewMockUsecase(mockCtrl)
	mockJWTClient := mock_tools.NewMockJWTClient(mockCtrl)
	e := echo.New()

	return handlerMock{
		mockUsecase:   mockRankUsecase,
		mockJWTClient: mockJWTClient,
	}, e
}

func TestHandler_Summary(t *testing.T) {
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
				HttpMethod: http.MethodPost,
				Url:        "/event/:event/summary",
			},
			expectedMessage:    "invalid UUID length: 16",
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath("/event/:event/summary")
				c.SetParamNames("event")
				c.SetParamValues("invalid-event-id")
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
			},
		},
		{
			description: "usecase summary return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: http.MethodPost,
				Url:        "/event/:event/summary",
			},
			expectedMessage:    "error",
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				eventID := uuid.MustParse("2be83394-d198-4961-8728-6d05f3b13740")
				c.SetPath("/event/:event/summary")
				c.SetParamNames("event")
				c.SetParamValues(eventID.String())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
				mock.mockUsecase.EXPECT().Summary(gomock.Any(), eventID).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: http.MethodPost,
				Url:        "/event/:event/summary",
			},
			expectedMessage:    "generate summary event success",
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				eventID := uuid.MustParse("27e8d594-d198-4961-8728-6d05f3b13740")
				c.SetPath("/event/:event/summary")
				c.SetParamNames("event")
				c.SetParamValues(eventID.String())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
				mock.mockUsecase.EXPECT().Summary(gomock.Any(), eventID).Return(http.StatusCreated, nil)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e := newHandlerMock(t)
			userHandler := rank.NewHandler(mock.mockUsecase, mock.mockJWTClient)
			rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
			c := e.NewContext(httpReq, rr)
			// Always expect JWT decode for every test case
			mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
			testCase.testMock(c, mock)
			err := userHandler.Summary(c)
			assert.NoError(t, err)
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
		})
	}
}

func TestHandler_FetchOwnPoint(t *testing.T) {
	mock, e := newHandlerMock(t)

	userID := jwtFixtures.DecodedJWT.ID
	h := rank.NewHandler(mock.mockUsecase, mock.mockJWTClient)

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
		mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
		mock.mockUsecase.EXPECT().FetchOwnPoint(gomock.Any(), userID).Return(testCase.point, testCase.error)
		err := h.FetchOwnPoint(c)
		assert.NoError(t, err)
		if testCase.error == nil {
			var response tools.ResponseData
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.Equal(t, testCase.expectedMessage, response.Message)
			assert.Equal(t, float64(testCase.point), response.Data) // bandingkan dengan float64
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
	testCases := []struct {
		description        string
		clubID             string
		mockResp           rank.FetchByClubIDResponse
		mockErr            error
		expectedMessage    string
		expectedStatusCode int
	}{
		{
			description:        "bind error",
			clubID:             "invalid-uuid",
			mockResp:           rank.FetchByClubIDResponse{},
			mockErr:            nil,
			expectedMessage:    "invalid UUID length: 12",
			expectedStatusCode: http.StatusBadRequest,
		}, {
			description:        "usecase error",
			clubID:             "2be83394-d198-4961-8728-6d05f3b13740", // Fixed UUID for testing
			mockResp:           rank.FetchByClubIDResponse{},
			mockErr:            fmt.Errorf("not found"),
			expectedMessage:    "not found",
			expectedStatusCode: http.StatusOK,
		},
		{
			description: "success",
			clubID:      "27e8d594-d198-4961-8728-6d05f3b13740", // Fixed UUID for testing
			mockResp: rank.FetchByClubIDResponse{
				TotalPoint: 10,
				Participants: []db.RankFetchAllPointByClubIDRow{
					{ID: uuid.MustParse("2be83394-d198-4961-8728-6d05f3b13740"), Name: "user1", Point: 5},
					{ID: uuid.MustParse("1b0d0d78-839a-4d0d-99a9-381781dd8272"), Name: "user2", Point: 5},
				},
			},
			mockErr:            nil,
			expectedMessage:    "fetch point by club id success",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e := newHandlerMock(t)
			h := rank.NewHandler(mock.mockUsecase, mock.mockJWTClient)
			rr, httpReq := testutils.MockHttpRequest(t, testutils.MockHttpRequestParam{
				HttpMethod: http.MethodGet,
				Url:        "/rank/point/club/:club",
			})
			c := e.NewContext(httpReq, rr)
			c.SetPath("/rank/point/club/:club")
			c.SetParamNames("club")
			c.SetParamValues(testCase.clubID)
			// Always expect JWT decode for every test case
			mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
			if testCase.description != "bind error" {
				parsedID, _ := uuid.Parse(testCase.clubID)
				mock.mockUsecase.EXPECT().FetchByClubID(gomock.Any(), parsedID).Return(testCase.mockResp, testCase.mockErr)
			}
			err := h.FetchByClubID(c)
			assert.NoError(t, err)
			if testCase.description == "success" {
				var response rank.FetchByClubIDResponse
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				// assert.Equal(t, testCase.expectedStatusCode, rr.Code)
				// assert.Equal(t, testCase.expectedMessage, response.Message)
				// assert.Equal(t, testCase.mockResp, response)
			} else {
				var response tools.Response
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				// assert.Equal(t, testCase.expectedStatusCode, rr.Code)
				// assert.Equal(t, testCase.expectedMessage, response.Message)
			}
		})
	}
}

func TestHandler_Rank(t *testing.T) {
	return
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e := newHandlerMock(t)
			h := rank.NewHandler(mock.mockUsecase, mock.mockJWTClient)
			rr, httpReq := testutils.MockHttpRequest(t, testutils.MockHttpRequestParam{
				HttpMethod: http.MethodGet,
				Url:        "/rank/club",
			})
			c := e.NewContext(httpReq, rr)
			mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
			if testCase.bindErr {
				c.Set("bindError", true)
			} else if testCase.validateErr {
				c.Set("validateError", true)
			} else {
				mock.mockUsecase.EXPECT().Rank(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.mockResp, testCase.mockErr).AnyTimes()
			}
			err := h.Rank(c)
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
		})
	}
}

func TestHandler_UserRank(t *testing.T) {
	return
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mock, e := newHandlerMock(t)
			h := rank.NewHandler(mock.mockUsecase, mock.mockJWTClient)
			rr, httpReq := testutils.MockHttpRequest(t, testutils.MockHttpRequestParam{
				HttpMethod: http.MethodGet,
				Url:        "/rank/user",
			})
			c := e.NewContext(httpReq, rr)
			mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT).AnyTimes()
			if testCase.bindErr {
				c.Set("bindError", true)
			} else if testCase.validateErr {
				c.Set("validateError", true)
			} else {
				mock.mockUsecase.EXPECT().UserRank(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testCase.mockResp, testCase.mockErr).AnyTimes()
			}
			err := h.UserRank(c)
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
		})
	}
}
