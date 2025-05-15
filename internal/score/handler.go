package score

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

	//Order
	e.POST("/event/:event/bracket/:bracket/score/order", handler.storeOrUpdateOrder, m.JWTMiddleware(), m.Middleware.EventPrivilegeWithoutReviewer, m.Middleware.IsEventTurnLocked, m.Middleware.EventManipulationRemarkOngoing)
	e.GET("/event/:event/bracket/:bracket/score/order", handler.fetchOneOrder, m.JWTMiddleware())

	//Single
	e.POST("/event/:event/bracket/:bracket/score/single", handler.storeOrUpdateSingle, m.JWTMiddleware(), m.Middleware.EventPrivilegeWithoutReviewer, m.Middleware.IsEventTurnLocked, m.Middleware.EventManipulationRemarkOngoing)
	e.GET("/event/:event/bracket/:bracket/score/single", handler.fetchOneSingle, m.JWTMiddleware())

	//All
	e.PATCH("/event/:event/class/:class/score/lock", handler.lock, m.JWTMiddleware(), m.Middleware.EventPrivilegeWithoutReviewer, m.Middleware.IsEventTurnLocked, m.Middleware.EventManipulationRemarkOngoing)
}

func (h handler) storeOrUpdateOrder(c echo.Context) error {
	var arg orderStoreOrUpdateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.storeOrUpdateOrder(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "store or update order scores success"})
}

func (h handler) fetchOneOrder(c echo.Context) error {
	var arg fetchOneParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	score, err := h.usecase.fetchOneOrder(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one order score success", Data: score})
}

func (h handler) lock(c echo.Context) error {
	var arg lockParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.lock(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "lock score success"})
}

func (h handler) storeOrUpdateSingle(c echo.Context) error {
	var arg singleStoreOrUpdateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.storeOrUpdateSingle(ctx, arg)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "store or update single scores success"})
}

func (h handler) fetchOneSingle(c echo.Context) error {
	var arg fetchOneParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	score, err := h.usecase.fetchOneSingle(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one single score success", Data: score})
}
