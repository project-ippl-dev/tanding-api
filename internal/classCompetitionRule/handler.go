package classCompetitionRule

import (
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	handler := handler{usecase: usecase}

	e.POST("/class/rules", handler.store, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/class/rules", handler.fetchAll, m.JWTMiddleware())
	e.GET("/class/rules/:rules", handler.fetchOne, m.JWTMiddleware())
	e.PUT("/class/rules/:rules", handler.update, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.DELETE("/class/rules/:rules", handler.delete, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
}

func (h handler) store(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.store(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "store class competition rules success"})
}

func (h handler) fetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchAll(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch class competition rules success", pagination))
}

func (h handler) fetchOne(c echo.Context) error {
	ruleIDParam := c.Param("rules")
	ruleID, err := strconv.ParseInt(ruleIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	rule, err := h.usecase.fetchOne(ctx, ruleID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one competition rules sucess", Data: rule})
}

func (h handler) update(c echo.Context) error {
	var req request
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
	if err := h.usecase.update(ctx, req, ruleID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update class competition rules success"})
}

func (h handler) delete(c echo.Context) error {
	ruleIDParam := c.Param("rules")
	ruleID, err := strconv.ParseInt(ruleIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.delete(ctx, ruleID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete class competition rules success"})
}
