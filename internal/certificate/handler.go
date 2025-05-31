package certificate

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{usecase: usecase}
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	handler := NewHandler(usecase)

	e.GET("/certificate/:certificate", handler.FetchOne, m.JWTMiddleware())
	e.GET("/certificate/user/:user", handler.FetchByUserID, m.JWTMiddleware())
	e.GET("/certificate/club/:club", handler.FetchByClubID, m.JWTMiddleware())
}

func (h Handler) FetchOne(c echo.Context) error {
	certificateIDParam := c.Param("certificate")
	certificateID, err := uuid.Parse(certificateIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	certificate, err := h.usecase.FetchOne(ctx, certificateID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one certificate success", Data: certificate})
}

func (h Handler) FetchByUserID(c echo.Context) error {
	userID := c.Param("user")
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchByUserID(ctx, page, pageSize, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch certificate by user id success", pagination))
}

func (h Handler) FetchByClubID(c echo.Context) error {
	clubIDParam := c.Param("club")
	clubID, err := uuid.Parse(clubIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchByClubID(ctx, page, pageSize, clubID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch certificate by club id success", pagination))
}
