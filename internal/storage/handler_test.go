package storage_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/storage"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_storage "github.com/project-ippl-dev/tanding-api/mocks/storage"
	"github.com/project-ippl-dev/tanding-api/testutils"
)

type handlerMock struct {
	mockUsecase *mock_storage.MockUsecase
}

func newHandlerMock(t *testing.T) (handlerMock, *echo.Echo) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorageUsecase := mock_storage.NewMockUsecase(mockCtrl)

	e := echo.New()

	return handlerMock{
		mockUsecase: mockStorageUsecase,
	}, e
}

func TestHandler_Upload(t *testing.T) {
	mock, e := newHandlerMock(t)

	storageHandler := storage.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodPost
	path := "/storage/upload"

	fileUsecaseMimeTypeError, fileUsecaseMimeTypeErrorContentType := testutils.GenerateMockFile(t, "test")
	fileValidateError, fileValidateErrorContentType := testutils.GenerateMockFile(t, "test")
	fileUploadError, fileUploadErrorContentType := testutils.GenerateMockFile(t, "test")
	fileUploadSuccess, fileUploadSuccessContentType := testutils.GenerateMockFile(t, "test")

	testCases := []struct {
		description        string
		req                testutils.MockHttpRequestParam
		expectedResponse   interface{}
		expectedErr        bool
		expectedStatusCode int
		testMock           func(c echo.Context, mock handlerMock)
	}{
		{
			description: "formFile return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   "request Content-Type isn't multipart/form-data",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase ReadMIMEType return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqRawBody: &fileUsecaseMimeTypeError,
				ReqHeaders: map[string]string{
					echo.HeaderContentType: fileUsecaseMimeTypeErrorContentType.FormDataContentType(),
				},
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().ReadMIMEType(gomock.Any(), gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			description: "validate return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqRawBody: &fileValidateError,
				ReqHeaders: map[string]string{
					echo.HeaderContentType: fileValidateErrorContentType.FormDataContentType(),
				},
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().ReadMIMEType(gomock.Any(), gomock.Any()).SetArg(1, storage.FileInformation{
					FileName: "test",
					Size:     6 * 1024 * 1024,
					MIMEType: "application/pdf",
					Ext:      "pdf",
					Dir:      "photo",
				}).
					Return(nil)
			},
		},
		{
			description: "usecase upload return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqRawBody: &fileUploadError,
				ReqHeaders: map[string]string{
					echo.HeaderContentType: fileUploadErrorContentType.FormDataContentType(),
				},
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().ReadMIMEType(gomock.Any(), gomock.Any()).SetArg(1, storage.FileInformation{
					FileName: "test",
					Size:     3 * 1024 * 1024,
					MIMEType: "application/pdf",
					Ext:      "pdf",
					Dir:      "photo",
				}).
					Return(nil)
				mock.mockUsecase.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("", fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqRawBody: &fileUploadSuccess,
				ReqHeaders: map[string]string{
					echo.HeaderContentType: fileUploadSuccessContentType.FormDataContentType(),
				},
			},
			expectedResponse:   tools.ResponseData{Message: "upload success", Data: "https://google.com"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().ReadMIMEType(gomock.Any(), gomock.Any()).SetArg(1, storage.FileInformation{
					FileName: "test",
					Size:     3 * 1024 * 1024,
					MIMEType: "application/pdf",
					Ext:      "pdf",
					Dir:      "photo",
				}).
					Return(nil)
				mock.mockUsecase.EXPECT().Upload(gomock.Any(), gomock.Any()).Return("https://google.com", nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := storageHandler.Upload(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse string
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_Base64(t *testing.T) {
	mock, e := newHandlerMock(t)

	storageHandler := storage.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodPost
	path := "/storage/upload/base64"

	validReq := storage.Base64Params{
		Data: "base-64-here",
		Dir:  "photo",
	}

	invalidReq := storage.Base64Params{
		Data: "base-64-here",
		Dir:  "",
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
				ReqBody:    "invalid-body",
			},
			expectedResponse:   "code=400, message=Unmarshal type error: expected=storage.Base64Params, got=string, field=, offset=14, internal=json: cannot unmarshal string into Go value of type storage.Base64Params",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
			},
		},
		{
			description: "usecase Base64ToFile return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Base64ToFile(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
		{
			description: "usecase Base64ToFile return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   "error",
			expectedErr:        true,
			expectedStatusCode: http.StatusBadRequest,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Base64ToFile(gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
		},
		{
			description: "validate return error",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    invalidReq,
			},
			expectedResponse:   "error validation",
			expectedErr:        true,
			expectedStatusCode: http.StatusUnprocessableEntity,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Base64ToFile(gomock.Any()).Return(testutils.GenerateDummyFile(1), nil)
			},
		},
		{
			description: "usecase Base64Upload return error",
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
				mock.mockUsecase.EXPECT().Base64ToFile(gomock.Any()).Return(testutils.GenerateDummyFile(1), nil)
				mock.mockUsecase.EXPECT().Base64Upload(gomock.Any()).Return("", fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
				ReqBody:    validReq,
			},
			expectedResponse:   tools.ResponseData{Message: "upload success", Data: "https://google.com"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				mock.mockUsecase.EXPECT().Base64ToFile(gomock.Any()).Return(testutils.GenerateDummyFile(1), nil)
				mock.mockUsecase.EXPECT().Base64Upload(gomock.Any()).Return("https://google.com", nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := storageHandler.Base64(c)
		assert.NoError(t, err)
		assert.Equal(t, testCase.expectedStatusCode, rr.Code)

		if testCase.expectedErr {
			var response tools.Response
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResponse, response.Message)
		} else {
			var response, expectedResponse string
			expectedResponseBytes, _ := json.Marshal(testCase.expectedResponse)
			_ = json.Unmarshal(expectedResponseBytes, &expectedResponse)
			_ = json.Unmarshal(rr.Body.Bytes(), &response)

			assert.Equal(t, expectedResponse, response)
		}
	}
}

func TestHandler_Delete(t *testing.T) {
	mock, e := newHandlerMock(t)

	storageHandler := storage.NewHandler(mock.mockUsecase)

	httpMethod := http.MethodDelete
	path := "/storage/delete/:dir/:file"

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
			expectedStatusCode: http.StatusInternalServerError,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("dir", "file")
				c.SetParamValues("photo", "https://google.com")
				mock.mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(http.StatusInternalServerError, fmt.Errorf("error"))
			},
		},
		{
			description: "success",
			req: testutils.MockHttpRequestParam{
				HttpMethod: httpMethod,
				Url:        path,
			},
			expectedResponse:   tools.Response{Message: "delete file success"},
			expectedErr:        false,
			expectedStatusCode: http.StatusOK,
			testMock: func(c echo.Context, mock handlerMock) {
				c.SetPath(path)
				c.SetParamNames("dir", "file")
				c.SetParamValues("photo", "https://google.com")
				mock.mockUsecase.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(http.StatusOK, nil)
			},
		},
	}

	for _, testCase := range testCases {
		rr, httpReq := testutils.MockHttpRequest(t, testCase.req)
		c := e.NewContext(httpReq, rr)
		testCase.testMock(c, mock)
		err := storageHandler.Delete(c)
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
