package classCompetitionRule

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	handler := NewHandler(usecase)

	e.POST("/class/rules", handler.Store, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/class/rules", handler.FetchAll, m.JWTMiddleware())
	e.GET("/class/rules/:rules", handler.FetchOne, m.JWTMiddleware())
	e.PUT("/class/rules/:rules", handler.Update, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.DELETE("/class/rules/:rules", handler.Delete, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
}

func (h Handler) Store(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Store(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "store class competition rules success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch class competition rules success", pagination))
}

func (h Handler) FetchOne(c echo.Context) error {
	ruleIDParam := c.Param("rules")
	ruleID, err := strconv.ParseInt(ruleIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	rule, err := h.usecase.FetchOne(ctx, ruleID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one competition rules success", Data: rule})
}

func (h Handler) Update(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	ruleIDParam := c.Param("rules")
	ruleID, err := strconv.ParseInt(ruleIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Update(ctx, req, ruleID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update class competition rules success"})
}

func (h Handler) Delete(c echo.Context) error {
	ruleIDParam := c.Param("rules")
	ruleID, err := strconv.ParseInt(ruleIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Delete(ctx, ruleID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete class competition rules success"})
}
