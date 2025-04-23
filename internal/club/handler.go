package club

import (
	"github.com/google/uuid"
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

	e.POST("/club", handler.store, m.JWTMiddleware())
	e.PUT("/club/:club", handler.update, m.JWTMiddleware())
	e.DELETE("/club/:club", handler.delete, m.JWTMiddleware())
}

func (h handler) store(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	clubID, err := h.usecase.store(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.ResponseData{Message: "store club success", Data: clubID})
}

func (h handler) update(c echo.Context) error {
	clubParam := c.Param("club")
	clubID, err := uuid.Parse(clubParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	var req request
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err = req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	if err = h.usecase.update(ctx, req, decoded.ID, clubID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update club success"})
}

func (h handler) delete(c echo.Context) error {
	clubParam := c.Param("club")
	clubID, err := uuid.Parse(clubParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	if err = h.usecase.delete(ctx, decoded.ID, clubID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete club success"})
}
