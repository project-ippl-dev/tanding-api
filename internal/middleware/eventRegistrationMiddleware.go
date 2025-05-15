package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

func (m Middleware) EventRegistrationOnlyClubOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		decoded := m.jwtClient.Decode(c)
		ctx := c.Request().Context()
		if _, err := m.repository.ClubFetchOwner(ctx, uuid.MustParse(decoded.ID)); err != nil {
			return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access, only user which had club can access"})
		}
		return next(c)
	}
}
