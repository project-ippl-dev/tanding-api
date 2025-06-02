package certificate_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/certificate"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_certificate "github.com/project-ippl-dev/tanding-api/mocks/certificate"
	certificateFixtures "github.com/project-ippl-dev/tanding-api/mocks/fixtures/certificate"
	"github.com/project-ippl-dev/tanding-api/testutils"
)

type handlerMock struct {
	mockUsecase *mock_certificate.MockUsecase
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCertificateUsecase := mock_certificate.NewMockUsecase(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase: mockCertificateUsecase,
	}, e
}

func TestHandler_FetchOne(t *testing.T) {
	mock, e := newHandlerMock(t)

	certificateHandler := certificate.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodGet
	path := "/certificate/:certificate"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "parse certificate id error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 22",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("certificate")
				c.SetParamValues("invalid-certificate-id")
			},
		},
		{
			description: "usecase FetchOne return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusNotFound,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("certificate")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(certificate.Response{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.ResponseData{Message: "fetch one certificate success", Data: certificateFixtures.CertificateResponse},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("certificate")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchOne(gomock.Any(), gomock.Any()).Return(certificateFixtures.CertificateResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := certificateHandler.FetchOne(c)
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

	certificateHandler := certificate.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodGet
	path := "/certificate/user/:user"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
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
				c.SetParamNames("user")
				c.SetParamValues("valid-user")
				mock.mockUsecase.EXPECT().FetchByUserID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch certificate by user id success", certificateFixtures.CertificateFetchResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("user")
				c.SetParamValues("valid-user")
				mock.mockUsecase.EXPECT().FetchByUserID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(certificateFixtures.CertificateFetchResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := certificateHandler.FetchByUserID(c)
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

func TestHandler_FetchByClubID(t *testing.T) {
	mock, e := newHandlerMock(t)

	certificateHandler := certificate.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodGet
	path := "/certificate/club/:club"

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "parse club id error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "invalid UUID length: 15",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("club")
				c.SetParamValues("invalid-club-id")
			},
		},
		{
			description: "usecase FetchByClubID return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("club")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchByClubID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(tools.Pagination{}, errors.New("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.PaginationGetResponse("fetch certificate by club id success", certificateFixtures.CertificateFetchResponse),
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("club")
				c.SetParamValues(uuid.NewString())
				mock.mockUsecase.EXPECT().FetchByClubID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(certificateFixtures.CertificateFetchResponse, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := certificateHandler.FetchByClubID(c)
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
