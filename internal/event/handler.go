package event

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
	"net/http"
)

type handler struct {
	usecase   Usecase
	jwtClient tools.JWTClient
}

func RegisterHandler(usecase Usecase, m middleware.Params, jwtClient tools.JWTClient, e *echo.Echo) {
	eventHandler := handler{
		usecase:   usecase,
		jwtClient: jwtClient,
	}

	//Admin Role
	e.GET("/event", eventHandler.fetchAll, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.DELETE("/event/:event", eventHandler.delete, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.PATCH("/event/:event/status", eventHandler.updateStatus, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.PATCH("/event/:event/remark", eventHandler.updateRemark, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)

	//Only Auth
	e.GET("/event/:event", eventHandler.fetchOne, m.JWTMiddleware())
	e.GET("/event/infinite", eventHandler.fetchInfinite)
	e.POST("/event", eventHandler.store, m.JWTMiddleware())

	//Competition Privilege
	e.GET("/event/own", eventHandler.fetchByUser, m.JWTMiddleware(), m.Middleware.HasEventPrivilege)
	e.PUT("/event/:event", eventHandler.update, m.JWTMiddleware(), m.Middleware.GrantCompetition)

	//Event Owner
	e.POST("/event/:event/committee", eventHandler.assign, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
	e.GET("/event/:event/committee", eventHandler.committeeFetchAll, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
	e.PATCH("/event/:event/committee/:committee", eventHandler.committeeUpdate, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
	e.DELETE("/event/:event/committee/:committee", eventHandler.committeeDelete, m.JWTMiddleware(), m.Middleware.EventPrivilegeOwner)
}

func (h handler) store(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, eventID, err := h.usecase.store(ctx, req, decoded)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: "store event success", Data: eventID})
}

func (h handler) fetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchAll(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all event success", pagination))
}

func (h handler) update(c echo.Context) error {
	var req request
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
	statusCode, err := h.usecase.update(ctx, req, eventID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update event success"})
}

func (h handler) delete(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.delete(ctx, eventID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete event success"})
}

func (h handler) fetchByUser(c echo.Context) error {
	decoded := h.jwtClient.Decode(c)
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchByUser(ctx, page, pageSize, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all by auth user", pagination))
}

func (h handler) fetchOne(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	event, err := h.usecase.fetchOne(ctx, eventID, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one event success", Data: event})
}

func (h handler) fetchInfinite(c echo.Context) error {
	var args fetchInfiniteQueryParams
	if err := c.Bind(&args); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	limit := tools.PaginationLimit(c)
	ctx := c.Request().Context()
	result, err := h.usecase.fetchInfinite(ctx, limit, args)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h handler) updateStatus(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	var req statusReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.updateStatus(ctx, req, eventID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update status event success"})
}

func (h handler) assign(c echo.Context) error {
	var req assignRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.assign(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "assign user success"})
}

func (h handler) committeeFetchAll(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	ctx := c.Request().Context()
	committees, err := h.usecase.committeeFetchAll(ctx, eventID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch event committee success", Data: committees})
}

func (h handler) committeeUpdate(c echo.Context) error {
	var arg updateCommitteeParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	statusCode, err := h.usecase.committeeUpdate(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update committee role success"})
}

func (h handler) committeeDelete(c echo.Context) error {
	var arg deleteCommitteeParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	statusCode, err := h.usecase.committeeDelete(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "delete committee role success"})
}

func (h handler) updateRemark(c echo.Context) error {
	var arg updateRemarkParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.usecase.updateRemark(ctx, arg); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update remark event success"})
}
