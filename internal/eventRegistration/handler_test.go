package eventRegistration_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/eventRegistration"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_eventRegistration "github.com/project-ippl-dev/tanding-api/mocks/eventRegistration"
	eventRegistrationFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/eventRegistration"
	jwtFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/tools/jwt"
	mock_tools "github.com/project-ippl-dev/tanding-api/mocks/tools"
	"github.com/project-ippl-dev/tanding-api/testutils"
)

type handlerMock struct {
	mockUsecase   *mock_eventRegistration.MockUsecase
	mockJWTClient *mock_tools.MockJWTClient
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEventRegistrationUsecase := mock_eventRegistration.NewMockUsecase(mockCtrl)
	mockJWTClient := mock_tools.NewMockJWTClient(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase:   mockEventRegistrationUsecase,
		mockJWTClient: mockJWTClient,
	}, e
}

func TestHandler_Register(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventRegistrationHandler := eventRegistration.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPost
	path := "/event/:event/register"

	validReq := eventRegistration.RegistrationRequest{
		EventID:      uuid.NewString(),
		ClassEventID: uuid.NewString(),
		ClubID:       uuid.NewString(),
		Members: []eventRegistration.MemberData{
			{
				UserID: uuid.New(),
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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventRegistration.RegistrationRequest, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventRegistration.RegistrationRequest",
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
				ReqBody:    eventRegistration.RegistrationRequest{},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase Register return error",
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
				mock.mockUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "register to specific event success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventRegistrationHandler.Register(c)
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

	eventRegistrationHandler := eventRegistration.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/:event/register"

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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventRegistration.FetchAllParams, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventRegistration.FetchAllParams",
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
			expectedResponse:   tools.PaginationGetResponse("fetch all event registration by event id success", eventRegistrationFixtures.EventRegistrationFetchAllResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(eventRegistrationFixtures.EventRegistrationFetchAllResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventRegistrationHandler.FetchAll(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse tools.PaginationResponse
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_Update(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventRegistrationHandler := eventRegistration.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPatch
	path := "/event/:event/register/:register"

	validReq := eventRegistration.UpdateRegistrationRequest{
		EventID:      uuid.NewString(),
		ClassEventID: uuid.NewString(),
		RegisterID:   uuid.New(),
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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventRegistration.UpdateRegistrationRequest, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventRegistration.UpdateRegistrationRequest",
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
				ReqBody:    eventRegistration.UpdateRegistrationRequest{},
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
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "update event registration success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventRegistrationHandler.Update(c)
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

func TestHandler_SetReject(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventRegistrationHandler := eventRegistration.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPatch
	path := "/event/:event/register/:register/rejected"

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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventRegistration.SetStatusRequest, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventRegistration.SetStatusRequest",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase SetReject return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "register")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().SetReject(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.Response{Message: "set registration status reject success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event", "register")
				c.SetParamValues(uuid.NewString(), uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().SetReject(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventRegistrationHandler.SetReject(c)
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

func TestHandler_FetchParticipant(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventRegistrationHandler := eventRegistration.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/:event/participant"

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
			description: "usecase FetchParticipant return error",
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
				mock.mockUsecase.EXPECT().FetchParticipant(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.ResponseData{Message: "fetch event participants success", Data: eventRegistrationFixtures.EventRegistrationFetchParticipantRows},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchParticipant(gomock.Any(), gomock.Any()).Return(eventRegistrationFixtures.EventRegistrationFetchParticipantRows, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventRegistrationHandler.FetchParticipant(c)
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
