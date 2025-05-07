package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Middleware struct {
	repository *db.Queries
}

type Params struct {
	Middleware Middleware
}

func (p Params) JWTMiddleware() echo.MiddlewareFunc {
	return tools.JWTMiddleware()
}

func InitMiddleware(repository *db.Queries) Middleware {
	return Middleware{repository: repository}
}

func (m Middleware) GrantCompetition(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		decoded := tools.JWTDecode(c)
		if _, err := m.repository.PrivilegeFetchOneByUserID(ctx, db.PrivilegeFetchOneByUserIDParams{
			ID:   uuid.MustParse(decoded.ID),
			Name: "competition",
			Type: db.PrivilegeTypeMain,
		}); err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) GrantExclusive(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		decoded := tools.JWTDecode(c)
		if _, err := m.repository.PrivilegeFetchOneByUserID(ctx, db.PrivilegeFetchOneByUserIDParams{
			ID:   uuid.MustParse(decoded.ID),
			Name: "exclusive",
			Type: db.PrivilegeTypeMain,
		}); err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) GrantOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		decoded := tools.JWTDecode(c)
		if _, err := m.repository.PrivilegeFetchOneByUserID(ctx, db.PrivilegeFetchOneByUserIDParams{
			ID:   uuid.MustParse(decoded.ID),
			Name: "owner",
			Type: db.PrivilegeTypeCompetition,
		}); err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) GrantAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		decoded := tools.JWTDecode(c)
		if _, err := m.repository.PrivilegeFetchOneByUserID(ctx, db.PrivilegeFetchOneByUserIDParams{
			ID:   uuid.MustParse(decoded.ID),
			Name: "admin",
			Type: db.PrivilegeTypeCompetition,
		}); err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) GrantReviewer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		decoded := tools.JWTDecode(c)
		if _, err := m.repository.PrivilegeFetchOneByUserID(ctx, db.PrivilegeFetchOneByUserIDParams{
			ID:   uuid.MustParse(decoded.ID),
			Name: "reviewer",
			Type: db.PrivilegeTypeCompetition,
		}); err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) GrantContributor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		decoded := tools.JWTDecode(c)
		if _, err := m.repository.PrivilegeFetchOneByUserID(ctx, db.PrivilegeFetchOneByUserIDParams{
			ID:   uuid.MustParse(decoded.ID),
			Name: "contributor",
			Type: db.PrivilegeTypeCompetition,
		}); err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}

func (m Middleware) RoleAdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		decoded := tools.JWTDecode(c)
		if decoded.RoleName != "admin" {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
		}
		return next(c)
	}
}
