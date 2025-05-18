package bracket_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/bracket"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_bracket "github.com/project-ippl-dev/tanding-api/mocks/bracket"
	"github.com/project-ippl-dev/tanding-api/mocks/fixtures/bracket"
	eventFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/event"
	eventRegistrationFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/eventRegistration"
	"github.com/project-ippl-dev/tanding-api/testutils"
	"github.com/project-ippl-dev/tanding-api/utils/pointer"
)

type handlerMock struct {
	mockUsecase *mock_bracket.MockUsecase
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBracketUsecase := mock_bracket.NewMockUsecase(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase: mockBracketUsecase,
	}, e
}

func TestHandler_Store(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	classEventID := bracketFixtures.ClassEventID

	httpMethod := http.MethodPost
	url := "/event/:event/class/:class/bracket"

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
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-event-id")
			},
		},
		{
			description: "usecase Store return error ",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "error",
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().Store(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "generate bracket success",
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().Store(gomock.Any(), gomock.Any()).Return(http.StatusCreated, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := bracketHandler.Store(c)
		assert.NoError(t, err)

		var response tools.Response
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		assert.Equal(t, testCase.expectedMessage, response.Message)
	}
}

func TestHandler_FetchOne(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	classEventID := bracketFixtures.ClassEventID

	httpMethod := http.MethodGet
	url := "/event/:event/class/:class/bracket"

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
				Url:        url,
			},
			expectedResponse:   "code=400, message=invalid UUID length: 22, internal=invalid UUID length: 22",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-event-id")
			},
		},
		{
			description: "usecase FetchOne return error ",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        true,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(bracket.FetchOneResponse{}, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   bracketFixtures.FetchOneResponseMatchTypeOrder,
			expectedStatusCode: http.StatusOK,
			expectedErr:        false,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(bracketFixtures.FetchOneResponseMatchTypeOrder, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := bracketHandler.FetchOne(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse bracket.FetchOneResponse
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_RoundDown(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	classEventID := bracketFixtures.ClassEventID

	httpMethod := http.MethodGet
	url := "/event/:event/class/:class/bracket/random"

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
				Url:        url,
			},
			expectedResponse:   "code=400, message=invalid UUID length: 22, internal=invalid UUID length: 22",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-event-id")
			},
		},
		{
			description: "usecase RoundDown return error ",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        true,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().RoundDown(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, bracket.RoundDownResponse{}, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   bracketFixtures.RoundDownResponseMatchTypeOrder,
			expectedStatusCode: http.StatusOK,
			expectedErr:        false,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().RoundDown(gomock.Any(), gomock.Any()).Return(http.StatusOK, bracketFixtures.RoundDownResponseMatchTypeOrder, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := bracketHandler.RoundDown(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse bracket.FetchOneResponse
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_UpdateOrderLock(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	classEventID := bracketFixtures.ClassEventID

	httpMethod := http.MethodPatch
	url := "/event/:event/class/:class/bracket/order/lock"

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
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-event-id")
			},
		},
		{
			description: "validation error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody: bracket.UpdateLockParams{
					ClassEventID: classEventID,
					Status:       pointer.ConvertToPointer(true),
					Participants: nil,
				},
			},
			expectedMessage:    "participants: cannot be blank.",
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
			},
		},
		{
			description: "usecase OrderLock return error ",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody: bracket.UpdateLockParams{
					ClassEventID: classEventID,
					Status:       pointer.ConvertToPointer(true),
					Participants: []bracket.ParticipantLockParams{
						{
							EventRegistrationID: eventRegistrationFixtures.EventRegistrationID.String(),
							Iteration:           1,
						},
					},
				},
			},
			expectedMessage:    "error",
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().OrderLock(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody: bracket.UpdateLockParams{
					ClassEventID: classEventID,
					Status:       pointer.ConvertToPointer(true),
					Participants: []bracket.ParticipantLockParams{
						{
							EventRegistrationID: eventRegistrationFixtures.EventRegistrationID.String(),
							Iteration:           1,
						},
					},
				},
			},
			expectedMessage:    "update lock status on specific bracket success",
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().OrderLock(gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := bracketHandler.UpdateOrderLock(c)
		assert.NoError(t, err)

		var response tools.Response
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		assert.Equal(t, testCase.expectedMessage, response.Message)
	}
}

func TestHandler_CancelBracket(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	classEventID := bracketFixtures.ClassEventID

	httpMethod := http.MethodPatch
	url := "/event/:event/class/:class/bracket/generate/canceled"

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
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-event-id")
			},
		},
		{
			description: "usecase CancelBracket return error ",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "error",
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().CancelBracket(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "update generate status on specific bracket success",
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().CancelBracket(gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := bracketHandler.CancelBracket(c)
		assert.NoError(t, err)

		var response tools.Response
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		assert.Equal(t, testCase.expectedMessage, response.Message)
	}
}

func TestHandler_UpdateSingleLock(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	classEventID := bracketFixtures.ClassEventID

	httpMethod := http.MethodPatch
	url := "/event/:event/class/:class/bracket/single/lock"

	validReq := bracket.UpdateSingleLockParams{
		EventID:      eventFixtures.EventID,
		ClassEventID: classEventID,
		Data:         bracketFixtures.BracketMatchIndexData,
		Status:       pointer.ConvertToPointer(true),
	}

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
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-event-id")
			},
		},
		{
			description: "usecase UpdateSingleLock return error ",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody:    validReq,
			},
			expectedMessage:    "error",
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().UpdateSingleLock(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody:    validReq,
			},
			expectedMessage:    "lock bracket for single elimination success",
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("class")
				c.SetParamValues(classEventID.String())
				mock.mockUsecase.EXPECT().UpdateSingleLock(gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := bracketHandler.UpdateSingleLock(c)
		assert.NoError(t, err)

		var response tools.Response
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		assert.Equal(t, testCase.expectedMessage, response.Message)
	}
}

func TestHandler_EventTurnLock(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	eventID := eventFixtures.EventID

	httpMethod := http.MethodPatch
	url := "/event/:event/turn/lock"

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
			expectedMessage:    "invalid UUID length: 16",
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("event")
				c.SetParamValues("invalid-event-id")
			},
		},
		{
			description: "usecase EventTurnLock return error ",
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
				mock.mockUsecase.EXPECT().EventTurnLock(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedMessage:    "generate event turn success",
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				c.SetParamNames("event")
				c.SetParamValues(eventID.String())
				mock.mockUsecase.EXPECT().EventTurnLock(gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := bracketHandler.EventTurnLock(c)
		assert.NoError(t, err)

		var response tools.Response
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)
		assert.Equal(t, testCase.expectedMessage, response.Message)
	}
}
