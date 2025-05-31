package score

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

	//Order
	e.POST("/event/:event/bracket/:bracket/score/order", handler.StoreOrUpdateOrder, m.JWTMiddleware(), m.Middleware.EventPrivilegeWithoutReviewer, m.Middleware.IsEventTurnLocked, m.Middleware.EventManipulationRemarkOngoing)
	e.GET("/event/:event/bracket/:bracket/score/order", handler.FetchOneOrder, m.JWTMiddleware())

	//Single
	e.POST("/event/:event/bracket/:bracket/score/single", handler.StoreOrUpdateSingle, m.JWTMiddleware(), m.Middleware.EventPrivilegeWithoutReviewer, m.Middleware.IsEventTurnLocked, m.Middleware.EventManipulationRemarkOngoing)
	e.GET("/event/:event/bracket/:bracket/score/single", handler.FetchOneSingle, m.JWTMiddleware())

	//All
	e.PATCH("/event/:event/class/:class/score/lock", handler.Lock, m.JWTMiddleware(), m.Middleware.EventPrivilegeWithoutReviewer, m.Middleware.IsEventTurnLocked, m.Middleware.EventManipulationRemarkOngoing)
}

func (h Handler) StoreOrUpdateOrder(c echo.Context) error {
	var arg OrderStoreOrUpdateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.StoreOrUpdateOrder(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "store or update order scores success"})
}

func (h Handler) FetchOneOrder(c echo.Context) error {
	var arg FetchOneParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	score, err := h.usecase.FetchOneOrder(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one order score success", Data: score})
}

func (h Handler) Lock(c echo.Context) error {
	var arg LockParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Lock(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "lock score success"})
}

func (h Handler) StoreOrUpdateSingle(c echo.Context) error {
	var arg SingleStoreOrUpdateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.StoreOrUpdateSingle(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "store or update single scores success"})
}

func (h Handler) FetchOneSingle(c echo.Context) error {
	var arg FetchOneParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	score, err := h.usecase.FetchOneSingle(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one single score success", Data: score})
}
