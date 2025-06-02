package sport_test

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/sport"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	sportFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/sport"
	mock_sport "github.com/project-ippl-dev/tanding-api/mocks/sport"
	"github.com/project-ippl-dev/tanding-api/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

type handlerMock struct {
	mockUsecase *mock_sport.MockUsecase
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockSportUsecase := mock_sport.NewMockUsecase(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase: mockSportUsecase,
	}, e
}

func TestHandler_Store(t *testing.T) {
	mock, e := newHandlerMock(t)

	sportHandler := sport.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodPost
	path := "/event/:event/register"

	validReq := sport.Request{
		Name:        "name",
		Description: "description",
		SportType:   db.SportTypeSport,
		Thumbnail:   "https://google.com",
	}

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "bind error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    "invalid-req-body",
			},
			expectedResponse:   "code=400, message=Unmarshal type error: expected=sport.Request, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type sport.Request",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "validate return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    sport.Request{},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase Store return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "store sport success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := sportHandler.Store(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse tools.Response
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_FetchAll(t *testing.T) {
	mock, e := newHandlerMock(t)

	sportHandler := sport.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodGet
	path := "/sport"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "bind error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    "invalid-req-body",
			},
			expectedResponse:   "code=400, message=Unmarshal type error: expected=sport.FetchAllQueryParams, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type sport.FetchAllQueryParams",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase FetchAll return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch all sport success", sportFixtures.SportFetchAllResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sportFixtures.SportFetchAllResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := sportHandler.FetchAll(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse tools.ResponseData
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_Update(t *testing.T) {
	mock, e := newHandlerMock(t)

	sportHandler := sport.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodPost
	path := "/sport/:sport"

	validReq := sport.Request{
		Name:        "name",
		Description: "description",
		SportType:   db.SportTypeSport,
		Thumbnail:   "https://google.com",
	}

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "bind error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    "invalid-req-body",
			},
			expectedResponse:   "code=400, message=Unmarshal type error: expected=sport.Request, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type sport.Request",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "validate return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    sport.Request{},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase Update return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusNotFound,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "update sport success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := sportHandler.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse tools.Response
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	mock, e := newHandlerMock(t)

	sportHandler := sport.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodDelete
	path := "/sport/:sport"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "usecase Delete return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusNotFound,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("sport")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.Response{Message: "delete sport success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("sport")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := sportHandler.Delete(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse tools.Response
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}
