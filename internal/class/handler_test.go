package class_test

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/class"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_class "github.com/project-ippl-dev/tanding-api/mocks/class"
	classFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/class"
	jwtFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/tools/jwt"
	mock_tools "github.com/project-ippl-dev/tanding-api/mocks/tools"
	"github.com/project-ippl-dev/tanding-api/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

type handlerMock struct {
	mockUsecase   *mock_class.MockUsecase
	mockJWTClient *mock_tools.MockJWTClient
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClassUsecase := mock_class.NewMockUsecase(mockCtrl)
	mockJWTClient := mock_tools.NewMockJWTClient(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase:   mockClassUsecase,
		mockJWTClient: mockJWTClient,
	}, e
}

func TestHandler_Store(t *testing.T) {
	mock, e := newHandlerMock(t)

	classHandler := class.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPost
	path := "/class"

	invalidReq := class.Request{
		SportID:     uuid.NewString(),
		Name:        "U-45M",
		ClassRuleID: 1,
		Type:        db.ClassTypeDefault,
		MatchType:   db.MatchTypeSingle,
	}

	validReq := class.Request{
		SportID:     uuid.NewString(),
		Name:        "U-45M",
		ClassRuleID: 1,
		Type:        db.ClassTypeCustom,
		MatchType:   db.MatchTypeSingle,
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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=class.Request, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type class.Request",
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
				ReqBody:    class.Request{},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "forbidden action",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    invalidReq,
			},
			expectedResponse:   "forbidden access",
			expectedErr:        true,
			expectedStatusCode: http.StatusForbidden,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
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
			expectedResponse:   tools.Response{Message: "store class success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusCreated,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classHandler.Store(c)
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

	classHandler := class.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/class"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
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
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch all class success", classFixtures.ClassFetchAllResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(classFixtures.ClassFetchAllResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classHandler.FetchAll(c)
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

func TestHandler_FetchBySportID(t *testing.T) {
	mock, e := newHandlerMock(t)

	classHandler := class.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	path := "/class/sport/:sport"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "invalid sport id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("sport")
				c.SetParamValues("invalid-sport-id")
			},
		},
		{
			description: "usecase FetchBySportID return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("sport")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchBySportID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch all based on sport id", classFixtures.ClassFetchBySportIDResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("sport")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchBySportID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(classFixtures.ClassFetchBySportIDResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classHandler.FetchBySportID(c)
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

	classHandler := class.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPut
	path := "/class/:class"

	invalidReq := class.Request{
		SportID:     uuid.NewString(),
		Name:        "U-45M",
		ClassRuleID: 1,
		Type:        db.ClassTypeDefault,
		MatchType:   db.MatchTypeSingle,
	}

	validReq := class.Request{
		SportID:     uuid.NewString(),
		Name:        "U-45M",
		ClassRuleID: 1,
		Type:        db.ClassTypeCustom,
		MatchType:   db.MatchTypeSingle,
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
			expectedResponse:   "code=400, message=Unmarshal type error: expected=class.Request, got=string, field=, offset=18, internal=json: cannot unmarshal string into Go value of type class.Request",
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
				ReqBody:    class.Request{},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "forbidden action",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    invalidReq,
			},
			expectedResponse:   "forbidden access",
			expectedErr:        true,
			expectedStatusCode: http.StatusForbidden,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
			},
		},
		{
			description: "invalid class id",
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
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-id")
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
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
				c.SetParamNames("class")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
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
			expectedResponse:   tools.Response{Message: "update specific class success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("class")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classHandler.Update(c)
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

	classHandler := class.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodDelete
	path := "/class/:class"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "invalid class id",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 16",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("class")
				c.SetParamValues("invalid-class-id")
			},
		},
		{
			description: "usecase Delete return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("class")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.Response{Message: "delete class success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("class")
				c.SetParamValues(uuid.NewString())
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := classHandler.Delete(c)
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
