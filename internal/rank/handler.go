package rank

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

	e.POST("/event/:event/summary", handler.summary, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin)
	e.GET("/rank/point/own", handler.fetchOwnPoint, m.JWTMiddleware())
	e.GET("/rank/point/club/:club", handler.fetchByClubID, m.JWTMiddleware())
	e.GET("/rank/club", handler.rank, m.JWTMiddleware())
	e.GET("/rank/user", handler.userRank, m.JWTMiddleware())
}

func (h handler) summary(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.summary(ctx, eventID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "generate summary event success"})
}

func (h handler) fetchOwnPoint(c echo.Context) error {
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	point, err := h.usecase.fetchOwnPoint(ctx, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch point by auth user success", Data: point})
}

func (h handler) fetchByClubID(c echo.Context) error {
	clubIDParam := c.Param("club")
	clubID, err := uuid.Parse(clubIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	result, err := h.usecase.fetchByClubID(ctx, clubID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch point by club id success", Data: result})
}

func (h handler) rank(c echo.Context) error {
	var arg rankParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.rank(ctx, page, pageSize, arg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all rank success", pagination))
}

func (h handler) userRank(c echo.Context) error {
	var arg rankParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.userRank(ctx, page, pageSize, arg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all user rank success", pagination))
}
