package bracket

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
	bracketHandler := NewHandler(usecase)

	e.POST("/event/:event/class/:class/bracket", bracketHandler.Store, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.GET("/event/:event/class/:class/bracket", bracketHandler.FetchOne)
	e.GET("/event/:event/class/:class/bracket/random", bracketHandler.RoundDown, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/class/:class/bracket/order/lock", bracketHandler.UpdateOrderLock, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/class/:class/bracket/single/lock", bracketHandler.UpdateSingleLock, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/class/:class/bracket/generate/canceled", bracketHandler.CancelBracket, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
	e.PATCH("/event/:event/turn/lock", bracketHandler.EventTurnLock, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin, m.Middleware.EventManipulationRemarkClosed)
}

func (h Handler) Store(c echo.Context) error {
	var arg GenerateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Store(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "generate bracket success"})
}

func (h Handler) FetchOne(c echo.Context) error {
	var arg GenerateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	result, err := h.usecase.FetchOne(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h Handler) RoundDown(c echo.Context) error {
	var arg GenerateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, response, err := h.usecase.RoundDown(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, response)
}

func (h Handler) UpdateOrderLock(c echo.Context) error {
	var arg UpdateLockParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.OrderLock(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update lock status on specific bracket success"})
}

func (h Handler) CancelBracket(c echo.Context) error {
	var arg UpdateGenerateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.CancelBracket(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update generate status on specific bracket success"})
}

func (h Handler) UpdateSingleLock(c echo.Context) error {
	var arg UpdateSingleLockParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.UpdateSingleLock(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "lock bracket for single elimination success"})
}

func (h Handler) EventTurnLock(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.EventTurnLock(ctx, eventID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "generate event turn success"})
}
