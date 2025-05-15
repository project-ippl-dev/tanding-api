package eventRegistration

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type handler struct {
	usecase   Usecase
	jwtClient tools.JWTClient
}

func RegisterHandler(usecase Usecase, m middleware.Params, jwtClient tools.JWTClient, e *echo.Echo) {
	eventRegistrationHandler := handler{
		usecase:   usecase,
		jwtClient: jwtClient,
	}

	e.GET("/event/:event/register", eventRegistrationHandler.fetchAll, m.JWTMiddleware())
	e.POST("/event/:event/register", eventRegistrationHandler.register, m.JWTMiddleware(), m.Middleware.EventManipulationRemarkOpen, m.Middleware.EventRegistrationOnlyClubOwner)
	e.PATCH("/event/:event/register/:register", eventRegistrationHandler.update, m.JWTMiddleware())
	e.PATCH("/event/:event/register/:register/rejected", eventRegistrationHandler.setReject, m.JWTMiddleware())
	e.GET("/event/:event/participant", eventRegistrationHandler.fetchParticipant, m.JWTMiddleware())
}

func (h handler) register(c echo.Context) error {
	var req registrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.register(ctx, req)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "register to specific event success"})
}

func (h handler) fetchAll(c echo.Context) error {
	var args fetchAllParams
	if err := c.Bind(&args); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ownStatus := c.QueryParam("own")
	if ownStatus == "1" {
		decoded := h.jwtClient.Decode(c)
		args.UserID = decoded.ID
	}
	ctx := c.Request().Context()

	pagination, err := h.usecase.fetchAll(ctx, args, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all event registration by event id success", pagination))
}

func (h handler) update(c echo.Context) error {
	var req updateRegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.update(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update event registration success"})
}

func (h handler) setReject(c echo.Context) error {
	var arg setStatusRequest
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	statusCode, err := h.usecase.setReject(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "set registration status reject success"})
}

func (h handler) fetchParticipant(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	classParticipants, err := h.usecase.fetchParticipant(ctx, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch event participants success", Data: classParticipants})
}
