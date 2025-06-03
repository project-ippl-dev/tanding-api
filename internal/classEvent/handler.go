package classEvent

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Handler struct {
	usecase   Usecase
	jwtClient tools.JWTClient
}

func RegisterHandler(usecase Usecase, jwtClient tools.JWTClient, m middleware.Params, e *echo.Echo) {
	classEventHandler := NewHandler(usecase, jwtClient)

	e.POST("/event/:event/class/assign", classEventHandler.Assign, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
	e.GET("/event/:event/class", classEventHandler.FetchAll)
	e.DELETE("/event/:event/class/:class", classEventHandler.Detach, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner, m.Middleware.ClassEventManipulation)
	e.PUT("/event/:event/class/:class", classEventHandler.Update, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner, m.Middleware.ClassEventManipulation)
}

func NewHandler(usecase Usecase, jwtClient tools.JWTClient) Handler {
	return Handler{
		usecase:   usecase,
		jwtClient: jwtClient,
	}
}

func (h Handler) Assign(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	if err := h.usecase.Assign(ctx, decoded.ID, req, eventID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "success assign class to event"})
}

func (h Handler) Detach(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	classIDParam := c.Param("class")
	classID, err := uuid.Parse(classIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	if err := h.usecase.Detach(ctx, decoded.ID, eventID, classID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "success detach class to event"})
}

func (h Handler) FetchAll(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	classes, err := h.usecase.FetchAll(ctx, decoded.ID, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch all class event success", Data: classes})
}

func (h Handler) Update(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	classIDParam := c.Param("class")
	classID, err := uuid.Parse(classIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	var req updateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	if err := h.usecase.Update(ctx, req, decoded.ID, eventID, classID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update price class success"})
}
