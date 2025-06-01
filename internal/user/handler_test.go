package user_test

import (
	"encoding/json"
	"fmt"
	userFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/user"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"github.com/project-ippl-dev/tanding-api/internal/user"
	dbFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/db"
	jwtFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/tools/jwt"
	mock_tools "github.com/project-ippl-dev/tanding-api/mocks/tools"
	mock_user "github.com/project-ippl-dev/tanding-api/mocks/user"
	"github.com/project-ippl-dev/tanding-api/testutils"
)

type handlerMock struct {
	mockUsecase   *mock_user.MockUsecase
	mockJWTClient *mock_tools.MockJWTClient
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserUsecase := mock_user.NewMockUsecase(mockCtrl)
	mockJWTClient := mock_tools.NewMockJWTClient(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase:   mockUserUsecase,
		mockJWTClient: mockJWTClient,
	}, e
}

func TestHandler_Search(t *testing.T) {
	mock, e := newHandlerMock(t)

	userHandler := user.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	invalidUrl := fmt.Sprintf("/user/search?keyword=%s&limit=%s", "invalid", "invalid")
	validUrl := fmt.Sprintf("/user/search?keyword=%s&limit=%d", "valid", 1)

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
				Url:        invalidUrl,
			},
			expectedResponse:   "code=400, message=strconv.ParseInt: parsing \"invalid\": invalid syntax, internal=strconv.ParseInt: parsing \"invalid\": invalid syntax",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(invalidUrl)
			},
		},
		{
			description: "usecase Search return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        validUrl,
			},
			expectedResponse:   "error",
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        true,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(validUrl)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        validUrl,
			},
			expectedResponse: tools.ResponseData{
				Message: "fetch user by search success",
				Data:    dbFixtures.UserFetchByKeywordRows,
			},
			expectedStatusCode: http.StatusOK,
			expectedErr:        false,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(validUrl)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(dbFixtures.UserFetchByKeywordRows, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := userHandler.Search(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse []db.UserFetchByKeywordRow
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_FetchOne(t *testing.T) {
	mock, e := newHandlerMock(t)

	userHandler := user.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	url := "/profile/:uuid/basic"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "usecase FetchOne return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(user.BasicInformationResponse{}, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse: tools.ResponseData{
				Message: "fetch one basic information success",
				Data:    userFixtures.UserFetchBasicInformationResponse,
			},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(userFixtures.UserFetchBasicInformationResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := userHandler.FetchOne(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse user.BasicInformationResponse
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_Update(t *testing.T) {
	mock, e := newHandlerMock(t)

	userHandler := user.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodPut
	url := "/profile/:uuid/basic"

	invalidReq := user.UpdateBasicInformationParams{
		UserID:         "user-id",
		Name:           "name",
		BornAt:         "born-at",
		IdentityNumber: "123456",
		Phone:          "12345678901",
		Photo:          "https://google.com",
		Gender:         "male",
		About:          "about",
	}

	validReq := user.UpdateBasicInformationParams{
		UserID:         "user-id",
		Name:           "name",
		BornAt:         "born-at",
		IdentityNumber: "1234567890123456",
		Phone:          "12345678901",
		Photo:          "https://google.com",
		Gender:         "male",
		About:          "about",
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
				Url:        url,
				ReqBody:    "invalid body",
			},
			expectedResponse:   "code=400, message=Unmarshal type error: expected=user.UpdateBasicInformationParams, got=string, field=, offset=14, internal=json: cannot unmarshal string into Go value of type user.UpdateBasicInformationParams",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
			},
		},
		{
			description: "validate error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody:    invalidReq,
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
			},
		},
		{
			description: "usecase Update return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody:    validReq,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.Response{Message: "update biodata success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockJWTClient.EXPECT().Decode(gomock.Any()).Return(jwtFixtures.DecodedJWT)
				mock.mockUsecase.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := userHandler.Update(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		}
	}
}

func TestHandler_FetchAll(t *testing.T) {
	mock, e := newHandlerMock(t)

	userHandler := user.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	url := "/user"

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
				Url:        url,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch all user success", userFixtures.UserFetchAllResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockUsecase.EXPECT().FetchAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(userFixtures.UserFetchAllResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := userHandler.FetchAll(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse tools.Pagination
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_FetchLastLogin(t *testing.T) {
	mock, e := newHandlerMock(t)

	userHandler := user.NewHandler(mock.mockUsecase, mock.mockJWTClient)

	httpMethod := http.MethodGet
	url := "/user/login"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "usecase FetchLastLogin return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockUsecase.EXPECT().FetchLastLogin(gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        url,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch last login success", userFixtures.UserFetchLastLoginResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(url)
				mock.mockUsecase.EXPECT().FetchLastLogin(gomock.Any(), gomock.Any(), gomock.Any()).Return(userFixtures.UserFetchLastLoginResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := userHandler.FetchLastLogin(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse tools.Pagination
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}
