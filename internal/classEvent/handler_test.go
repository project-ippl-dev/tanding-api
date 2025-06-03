package classEvent_test

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/classEvent"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_classEvent "github.com/project-ippl-dev/tanding-api/mocks/classEvent"
	dbFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/db"
	jwtFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/tools/jwt"
	mock_tools "github.com/project-ippl-dev/tanding-api/mocks/tools"
	"github.com/project-ippl-dev/tanding-api/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

type handlerMock struct {
	mockUsecase   *mock_classEvent.MockUsecase
	mockJWTClient *mock_tools.MockJWTClient
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClassEventUsecase := mock_classEvent.NewMockUsecase(mockCtrl)
	mockJWTClient := mock_tools.NewMockJWTClient(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase:   mockClassEventUsecase,
		mockJWTClient: mockJWTClient,
	}, e
}

func TestHandler_Assign(t *testing.T) {
	mock, e := newHandlerMock(t)

	classEventHandler := classEvent.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPost
	path := "/event/:event/class/assign"

	validReq := classEvent.Request{
		Data: []classEvent.RequestData{
			{
				ClassID: uuid.NewString(),
				Price:   0,
			},
		},
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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=classEvent.Request, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type classEvent.Request",
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
				ReqBody:    classEvent.Request{},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "invalid event id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues("invalid-event-id")
			},
		},
		{
			description: "usecase Assign return error",
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
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Assign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "success assign class to event"},
			expectedErr:        false,
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Assign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classEventHandler.Assign(c)
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

func TestHandler_Detach(t *testing.T) {
	mock, e := newHandlerMock(t)

	classEventHandler := classEvent.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodDelete
	path := "/event/:event/class/:class"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "invalid event id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues("invalid-event-id", "invalid-class-id")
			},
		},
		{
			description: "invalid class event id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), "invalid-class-id")
			},
		},
		{
			description: "usecase Detach return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusNotFound,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Detach(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.Response{Message: "success detach class to event"},
			expectedErr:        false,
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Detach(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classEventHandler.Detach(c)
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

	classEventHandler := classEvent.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/:event/class"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "invalid event id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues("invalid-event-id")
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
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.ResponseData{Message: "fetch all class event success", Data: dbFixtures.ClassEventFetchAllRows},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(dbFixtures.ClassEventFetchAllRows, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classEventHandler.FetchAll(c)
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

	classEventHandler := classEvent.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPut
	path := "/event/:event/class/:class"

	validReq := classEvent.UpdateReq{
		Price: 1000,
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
			description: "invalid event id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues("invalid-event-id", "invalid-class-id")
			},
		},
		{
			description: "invalid class event id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), "invalid-class-id")
			},
		},
		{
			description: "bind error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    "invalid-req-body",
			},
			expectedResponse:   "code=400, message=Unmarshal type error: expected=classEvent.UpdateReq, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type classEvent.UpdateReq",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
			},
		},
		{
			description: "validate return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    classEvent.UpdateReq{Price: -1000},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
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
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "update price class success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "class")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classEventHandler.Update(c)
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
