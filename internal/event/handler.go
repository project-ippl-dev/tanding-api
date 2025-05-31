package event

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"net/http"
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
	eventHandler := NewHandler(usecase, jwtClient)

	//Admin Role
	e.GET("/event", eventHandler.FetchAll, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.DELETE("/event/:event", eventHandler.Delete, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.PATCH("/event/:event/status", eventHandler.UpdateStatus, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.PATCH("/event/:event/remark", eventHandler.UpdateRemark, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)

	//Only Auth
	e.GET("/event/:event", eventHandler.FetchOne, m.JWTMiddleware())
	e.GET("/event/infinite", eventHandler.FetchInfinite)
	e.POST("/event", eventHandler.Store, m.JWTMiddleware())

	//Competition Privilege
	e.GET("/event/own", eventHandler.FetchByUser, m.JWTMiddleware(), m.Middleware.HasEventPrivilege)
	e.PUT("/event/:event", eventHandler.Update, m.JWTMiddleware(), m.Middleware.GrantCompetition)

	//Event Owner
	e.POST("/event/:event/committee", eventHandler.Assign, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
	e.GET("/event/:event/committee", eventHandler.CommitteeFetchAll, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
	e.PATCH("/event/:event/committee/:committee", eventHandler.CommitteeUpdate, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
	e.DELETE("/event/:event/committee/:committee", eventHandler.CommitteeDelete, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
}

func (h Handler) Store(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, eventID, err := h.usecase.Store(ctx, req, decoded)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: "store event success", Data: eventID})
}

func (h Handler) FetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all event success", pagination))
}

func (h Handler) Update(c echo.Context) error {
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
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Update(ctx, req, eventID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update event success"})
}

func (h Handler) Delete(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Delete(ctx, eventID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete event success"})
}

func (h Handler) FetchByUser(c echo.Context) error {
	decoded := h.jwtClient.Decode(c)
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchByUser(ctx, page, pageSize, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all by auth user", pagination))
}

func (h Handler) FetchOne(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	event, err := h.usecase.FetchOne(ctx, eventID, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one event success", Data: event})
}

func (h Handler) FetchInfinite(c echo.Context) error {
	var args FetchInfiniteQueryParams
	if err := c.Bind(&args); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	limit := tools.PaginationLimit(c)
	ctx := c.Request().Context()
	result, err := h.usecase.FetchInfinite(ctx, limit, args)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h Handler) UpdateStatus(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	var req StatusReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.UpdateStatus(ctx, req, eventID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update status event success"})
}

func (h Handler) Assign(c echo.Context) error {
	var req AssignRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Assign(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "assign user success"})
}

func (h Handler) CommitteeFetchAll(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	ctx := c.Request().Context()
	committees, err := h.usecase.CommitteeFetchAll(ctx, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch event committee success", Data: committees})
}

func (h Handler) CommitteeUpdate(c echo.Context) error {
	var arg UpdateCommitteeParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	statusCode, err := h.usecase.CommitteeUpdate(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update committee role success"})
}

func (h Handler) CommitteeDelete(c echo.Context) error {
	var arg DeleteCommitteeParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	statusCode, err := h.usecase.CommitteeDelete(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "delete committee role success"})
}

func (h Handler) UpdateRemark(c echo.Context) error {
	var arg UpdateRemarkParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.usecase.UpdateRemark(ctx, arg); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update remark event success"})
}
