package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

func (m Middleware) ClassEventManipulation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		event, err := m.repository.EventCheckOne(ctx, eventID)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
		}
		if event.Remark != db.RemarkTypeUnconfirmed {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "can't update class event with event remark beside unconfirmed"})
		}
		return next(c)
	}
}

func (m Middleware) EventPrivilegeOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		decoded := m.jwtClient.Decode(c)
		privilege, err := m.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
			EventID: eventID,
			UserID:  uuid.MustParse(decoded.ID),
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, tools.Response{Message: "unauthorized"})
		}
		if privilege.Role != db.EventRoleOwner {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) EventPrivilegeWithoutContributor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		decoded := m.jwtClient.Decode(c)
		privilege, err := m.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
			EventID: eventID,
			UserID:  uuid.MustParse(decoded.ID),
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, tools.Response{Message: "unauthorized"})
		}
		if privilege.Role == db.EventRoleContributor {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) EventPrivilegeWithoutReviewer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		decoded := m.jwtClient.Decode(c)
		privilege, err := m.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
			EventID: eventID,
			UserID:  uuid.MustParse(decoded.ID),
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, tools.Response{Message: "unauthorized"})
		}
		if privilege.Role == db.EventRoleReviewer {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) EventPrivilegeAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		if eventIDParam == "" {
			eventIDParam = c.QueryParam("event_id")
		}
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		decoded := m.jwtClient.Decode(c)
		if decoded.RoleName == db.RoleAdmin {
			return next(c)
		}
		privilege, err := m.repository.EventPrivilegeFetchOne(ctx, db.EventPrivilegeFetchOneParams{
			EventID: eventID,
			UserID:  uuid.MustParse(decoded.ID),
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, tools.Response{Message: "unauthorized"})
		}
		if privilege.Role == db.EventRoleOwner || privilege.Role == db.EventRoleAdmin {
			return next(c)
		}
		return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
	}
}

func (m Middleware) HasEventPrivilege(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		decoded := m.jwtClient.Decode(c)
		events, err := m.repository.EventPrivilegeFetchByUserID(ctx, uuid.MustParse(decoded.ID))
		if err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		if len(events) == 0 {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) EventManipulationRemarkOngoing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		event, err := m.repository.EventCheckOne(ctx, eventID)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
		}
		if event.Remark != db.RemarkTypeOngoing {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "can't manipulation if event remark is not ongoing"})
		}
		return next(c)
	}
}

func (m Middleware) EventManipulationRemarkClosed(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		event, err := m.repository.EventCheckOne(ctx, eventID)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
		}
		if event.Remark != db.RemarkTypeClosed {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "can't manipulation if event remark is not closed"})
		}
		return next(c)
	}
}

func (m Middleware) EventManipulationRemarkOpen(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		event, err := m.repository.EventCheckOne(ctx, eventID)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
		}
		if event.Remark != db.RemarkTypeOpen {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "can't manipulation if event remark is not open"})
		}
		return next(c)
	}
}

func (m Middleware) IsEventTurnLocked(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		eventIDParam := c.Param("event")
		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
		}
		ctx := c.Request().Context()
		event, err := m.repository.EventCheckOne(ctx, eventID)
		if err != nil {
			return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
		}
		if !event.IsGenerate {
			return c.JSON(http.StatusUnprocessableEntity, tools.Response{Message: "lock and generate event turn first"})
		}
		return next(c)
	}
}
