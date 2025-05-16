package bracket_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/project-ippl-dev/tanding-api/internal/bracket"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	mock_bracket "github.com/project-ippl-dev/tanding-api/mocks/bracket"
	"github.com/project-ippl-dev/tanding-api/testutils"
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

func TestHandler_store(t *testing.T) {
	mock, e := newHandlerMock(t)

	bracketHandler := bracket.NewHandler(mock.mockUsecase)

	classEventID := uuid.New()

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
			description: "usecase store return error ",
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
