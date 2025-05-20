package certificate

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	handler := handler{usecase: usecase}

	e.GET("/certificate/:certificate", handler.fetchOne, m.JWTMiddleware())
	e.GET("/certificate/user/:user", handler.fetchByUserID, m.JWTMiddleware())
	e.GET("/certificate/club/:club", handler.fetchByClubID, m.JWTMiddleware())
}

func (h handler) fetchOne(c echo.Context) error {
	certificateIDParam := c.Param("certificate")
	certificateID, err := uuid.Parse(certificateIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	certificate, err := h.usecase.fetchOne(ctx, certificateID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one certificate success", Data: certificate})
}

func (h handler) fetchByUserID(c echo.Context) error {
	userID := c.Param("user")
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchByUserID(ctx, page, pageSize, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch certificate by user id success", pagination))
}

func (h handler) fetchByClubID(c echo.Context) error {
	clubIDParam := c.Param("club")
	clubID, err := uuid.Parse(clubIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchByClubID(ctx, page, pageSize, clubID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch certificate by club id success", pagination))
}
