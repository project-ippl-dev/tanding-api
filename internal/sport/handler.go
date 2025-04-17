package sport

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	handler := handler{usecase: usecase}

	e.POST("/sport", handler.store, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/sport", handler.fetchAll)
	e.PUT("/sport/:sport", handler.update, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.DELETE("/sport/:sport", handler.delete, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)

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
	return c.JSON(http.StatusCreated, tools.Response{Message: "store sport success"})
}

func (h handler) fetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	var queries fetchAllQueryParams
	if err := c.Bind(&queries); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchAll(ctx, page, pageSize, queries)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.PaginationGetResponse("fetch all sport success", pagination))
}

func (h handler) update(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	sportID := c.Param("sport")
	ctx := c.Request().Context()
	if err := h.usecase.update(ctx, req, sportID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update sport success"})
}

func (h handler) delete(c echo.Context) error {
	sportID := c.Param("sport")
	ctx := c.Request().Context()
	if err := h.usecase.delete(ctx, sportID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete sport success"})
}
