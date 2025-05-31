package sport

import (
	"net/http"

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

	e.POST("/sport", handler.Store, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/sport", handler.FetchAll)
	e.PUT("/sport/:sport", handler.Update, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.DELETE("/sport/:sport", handler.Delete, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)

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
	return c.JSON(http.StatusCreated, tools.Response{Message: "store sport success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	var queries FetchAllQueryParams
	if err := c.Bind(&queries); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize, queries)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.PaginationGetResponse("fetch all sport success", pagination))
}

func (h Handler) Update(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	sportID := c.Param("sport")
	ctx := c.Request().Context()
	if err := h.usecase.Update(ctx, req, sportID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update sport success"})
}

func (h Handler) Delete(c echo.Context) error {
	sportID := c.Param("sport")
	ctx := c.Request().Context()
	if err := h.usecase.Delete(ctx, sportID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete sport success"})
}
