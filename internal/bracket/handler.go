package bracket

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

	e.POST("/event/:event/class/:class/bracket", handler.store, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.GET("/event/:event/class/:class/bracket", handler.fetchOne)
	e.GET("/event/:event/class/:class/bracket/random", handler.roundDown, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/class/:class/bracket/order/lock", handler.updateOrderLock, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/class/:class/bracket/single/lock", handler.updateSingleLock, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/class/:class/bracket/generate/canceled", handler.cancelBracket, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/turn/lock", handler.eventTurnLock, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
}

func (h handler) store(c echo.Context) error {
	var arg generateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.store(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "generate bracket success"})
}

func (h handler) fetchOne(c echo.Context) error {
	var arg generateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	result, err := h.usecase.fetchOne(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h handler) roundDown(c echo.Context) error {
	var arg generateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, response, err := h.usecase.roundDown(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, response)
}

func (h handler) updateOrderLock(c echo.Context) error {
	var arg updateLockParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.orderLock(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update lock status on specific bracket success"})
}

func (h handler) cancelBracket(c echo.Context) error {
	var arg updateGenerateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.cancelBracket(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update generate status on specific bracket success"})
}

func (h handler) updateSingleLock(c echo.Context) error {
	var arg updateSingleLockParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.updateSingleLock(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "lock bracket for single elimination success"})
}

func (h handler) eventTurnLock(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.eventTurnLock(ctx, eventID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "generate event turn success"})
}
