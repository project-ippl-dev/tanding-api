package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/db"
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
	userHandler := NewHandler(usecase, jwtClient)

	e.GET("/user/search", userHandler.Search, m.JWTMiddleware())
	e.GET("/profile/:uuid/basic", userHandler.FetchOne, m.JWTMiddleware())
	e.PUT("/profile/:uuid/basic", userHandler.Update, m.JWTMiddleware())

	e.GET("/user", userHandler.FetchAll, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/user/login", userHandler.FetchLastLogin, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
}

func (h Handler) Search(c echo.Context) error {
	var args SearchParams
	if err := c.Bind(&args); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if args.Limit == 0 {
		args.Limit = 10
	}
	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	users, err := h.usecase.Search(ctx, args, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch user by search success", Data: users})
}

func (h Handler) FetchOne(c echo.Context) error {
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()

	userID := c.Param("uuid")
	if decoded.RoleName != db.RoleAdmin || userID == "own" {
		userID = decoded.ID
	}

	basic, err := h.usecase.FetchOne(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch one basic information success", Data: basic})
}

func (h Handler) Update(c echo.Context) error {
	var arg UpdateBasicInformationParams
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

	if err := h.usecase.Update(ctx, arg, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update biodata success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all user success", pagination))
}

func (h Handler) FetchLastLogin(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchLastLogin(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch last login success", pagination))
}
