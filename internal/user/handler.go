package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type handler struct {
	usecase   Usecase
	jwtClient tools.JWTClient
}

func RegisterHandler(usecase Usecase, m middleware.Params, jwtClient tools.JWTClient, e *echo.Echo) {
	userHandler := handler{
		usecase:   usecase,
		jwtClient: jwtClient,
	}

	e.GET("/user/search", userHandler.search, m.JWTMiddleware())
	e.GET("/profile/:uuid/basic", userHandler.fetchOne, m.JWTMiddleware())
	e.PUT("/profile/:uuid/basic", userHandler.update, m.JWTMiddleware())

	e.GET("/user", userHandler.fetchAll, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/user/login", userHandler.fetchLastLogin, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
}

func (h handler) search(c echo.Context) error {
	var args searchParams
	if err := c.Bind(&args); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if args.Limit == 0 {
		args.Limit = 10
	}
	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	users, err := h.usecase.search(ctx, args, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch user by search success", Data: users})
}

func (h handler) fetchOne(c echo.Context) error {
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()

	userID := c.Param("uuid")
	if decoded.RoleName != db.RoleAdmin || userID == "own" {
		userID = decoded.ID
	}

	basic, err := h.usecase.fetchOne(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one basic information success", Data: basic})
}

func (h handler) update(c echo.Context) error {
	var arg updateBasicInformationParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()

	userID := c.Param("uuid")
	if decoded.RoleName != db.RoleAdmin || userID == "own" {
		userID = decoded.ID
	}

	if err := h.usecase.update(ctx, arg, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update biodata success"})
}

func (h handler) fetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchAll(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all user success", pagination))
}

func (h handler) fetchLastLogin(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchLastLogin(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch last login success", pagination))
}
