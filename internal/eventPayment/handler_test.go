package eventPayment_test

import (
	"encoding/json"
	"errors"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	eventPaymentFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/eventPayment"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/eventPayment"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_eventPayment "github.com/project-ippl-dev/tanding-api/mocks/eventPayment"
	jwtFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/tools/jwt"
	mock_tools "github.com/project-ippl-dev/tanding-api/mocks/tools"
	"github.com/project-ippl-dev/tanding-api/testutils"
)

type handlerMock struct {
	mockUsecase   *mock_eventPayment.MockUsecase
	mockJWTClient *mock_tools.MockJWTClient
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockEventPaymentUsecase := mock_eventPayment.NewMockUsecase(mockCtrl)
	mockJWTClient := mock_tools.NewMockJWTClient(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase:   mockEventPaymentUsecase,
		mockJWTClient: mockJWTClient,
	}, e
}

func TestHandler_Store(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPost
	path := "/event/:event/payment"

	validReq := eventPayment.Request{
		EventID:       uuid.New(),
		UniqueNumber:  123,
		Link:          "https://google.com",
		Total:         1,
		Registrations: []uuid.UUID{uuid.New()},
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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventPayment.Request, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventPayment.Request",
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
				ReqBody:    eventPayment.Request{},
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
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "store event payment success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Store(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.Store(c)
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

func TestHandler_FetchAll(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/payment"

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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventPayment.FetchAllParams, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventPayment.FetchAllParams",
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
			expectedResponse:   tools.PaginationGetResponse("fetch all payment receipt success", eventPaymentFixtures.EventPaymentFetchResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(eventPaymentFixtures.EventPaymentFetchResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.FetchAll(c)
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

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPatch
	path := "/event/:event/payment/:payment"

	validReq := eventPayment.UpdateParams{
		EventID:   uuid.New(),
		PaymentID: uuid.New(),
		Status:    db.EventReceiptStatusApproved,
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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventPayment.UpdateParams, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventPayment.UpdateParams",
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
				ReqBody:    eventPayment.UpdateParams{},
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
			expectedResponse:   tools.Response{Message: "update payment for admin success"},
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
		err := eventPaymentHandler.Update(c)
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

func TestHandler_FetchByUserID(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/:event/payment/club"

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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventPayment.FetchByUserIDParams, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventPayment.FetchByUserIDParams",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase FetchByUserID return error",
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
				mock.mockUsecase.EXPECT().FetchByUserID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch by user id success", eventPaymentFixtures.EventPaymentFetchResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().FetchByUserID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, eventPaymentFixtures.EventPaymentFetchResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.FetchByUserID(c)
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

func TestHandler_FetchByEventPrivilege(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/:event/payment"

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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventPayment.FetchByEventPrivilegeParams, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventPayment.FetchByEventPrivilegeParams",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "bind error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "event param is required",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
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
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch all payment receipt in specific event success", eventPaymentFixtures.EventPaymentFetchResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(eventPaymentFixtures.EventPaymentFetchResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.FetchByEventPrivilege(c)
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

func TestHandler_Cart(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/cart"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "usecase Cart return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Cart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch cart by specific user success", eventPaymentFixtures.EventPaymentCartResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Cart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(eventPaymentFixtures.EventPaymentCartResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.Cart(c)
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

func TestHandler_CartDetail(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/:event/cart/detail"

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
			description: "usecase CartDetail return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusNotFound,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().CartDetail(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventPayment.CartDetailResponse{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.ResponseData{Message: "fetch cart detail success", Data: eventPaymentFixtures.EventPaymentCartDetailResponse},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("event")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().CartDetail(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventPaymentFixtures.EventPaymentCartDetailResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.CartDetail(c)
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

func TestHandler_Detail(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/payment/:payment/detail"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "invalid payment id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 18",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("payment")
				c.SetParamValues("invalid-payment-id")
			},
		},
		{
			description: "usecase Detail return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusNotFound,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("payment")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().Detail(gomock.Any(), gomock.Any()).Return(eventPayment.DetailPaymentResponse{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.ResponseData{Message: "fetch detail payment success", Data: eventPaymentFixtures.DetailPaymentResponse},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("payment")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().Detail(gomock.Any(), gomock.Any()).Return(eventPaymentFixtures.DetailPaymentResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.Detail(c)
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

func TestHandler_Summary(t *testing.T) {
	mock, e := newHandlerMock(t)

	eventPaymentHandler := eventPayment.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/event/:event/payment/summary"

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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=eventPayment.SummaryParams, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type eventPayment.SummaryParams",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase Summary return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusNotFound,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Summary(gomock.Any(), gomock.Any()).Return(eventPayment.SummaryResponse{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.ResponseData{Message: "fetch payment summary success", Data: eventPaymentFixtures.SummaryResponse},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Summary(gomock.Any(), gomock.Any()).Return(eventPaymentFixtures.SummaryResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := eventPaymentHandler.Summary(c)
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
