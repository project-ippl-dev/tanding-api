package eventRegistration

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

func NewHandler(usecase Usecase, jwtClient tools.JWTClient) Handler {
	return Handler{
		usecase:   usecase,
		jwtClient: jwtClient,
	}
}

func RegisterHandler(usecase Usecase, m middleware.Params, jwtClient tools.JWTClient, e *echo.Echo) {
	eventRegistrationHandler := NewHandler(usecase, jwtClient)

	e.GET("/event/:event/register", eventRegistrationHandler.FetchAll, m.JWTMiddleware())
	e.POST("/event/:event/register", eventRegistrationHandler.Register, m.JWTMiddleware(), m.Middleware.EventManipulationRemarkOpen, m.Middleware.EventRegistrationOnlyClubOwner)
	e.PATCH("/event/:event/register/:register", eventRegistrationHandler.Update, m.JWTMiddleware())
	e.PATCH("/event/:event/register/:register/rejected", eventRegistrationHandler.SetReject, m.JWTMiddleware())
	e.GET("/event/:event/participant", eventRegistrationHandler.FetchParticipant, m.JWTMiddleware())
}

func (h Handler) Register(c echo.Context) error {
	var req RegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Register(ctx, req)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "register to specific event success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	var args FetchAllParams
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

	pagination, err := h.usecase.FetchAll(ctx, args, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all event registration by event id success", pagination))
}

func (h Handler) Update(c echo.Context) error {
	var req UpdateRegistrationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Update(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update event registration success"})
}

func (h Handler) SetReject(c echo.Context) error {
	var arg SetStatusRequest
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	statusCode, err := h.usecase.SetReject(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "set registration status reject success"})
}

func (h Handler) FetchParticipant(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	classParticipants, err := h.usecase.FetchParticipant(ctx, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch event participants success", Data: classParticipants})
}
