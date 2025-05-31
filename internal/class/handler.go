package class

import (
	"net/http"

	"github.com/google/uuid"
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
	classHandler := NewHandler(usecase, jwtClient)

	e.POST("/class", classHandler.Store, m.JWTMiddleware())
	e.GET("/class", classHandler.FetchAll, m.JWTMiddleware())
	e.GET("/class/sport/:sport", classHandler.FetchBySportID, m.JWTMiddleware())
	e.PUT("/class/:class", classHandler.Update, m.JWTMiddleware())
	e.DELETE("/class/:class", classHandler.Delete, m.JWTMiddleware())
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
	if decoded.RoleName == "user" && req.Type != db.ClassTypeCustom {
		return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Store(ctx, req); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "store class success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all class success", pagination))
}

func (h Handler) FetchBySportID(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	sportIDParam := c.Param("sport")
	sportID, err := uuid.Parse(sportIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchBySportID(ctx, page, pageSize, sportID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all based on sport id", pagination))
}

func (h Handler) Update(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	if decoded.RoleName == "user" && req.Type != db.ClassTypeCustom {
		return c.JSON(http.StatusForbidden, tools.Response{Message: "forbidden access"})
	}

	classIDParam := c.Param("class")
	classID, err := uuid.Parse(classIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Update(ctx, req, classID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update specific class success"})
}

func (h Handler) Delete(c echo.Context) error {
	classIDParam := c.Param("class")
	classID, err := uuid.Parse(classIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	decoded := h.jwtClient.Decode(c)
	statusCode, err := h.usecase.Delete(ctx, decoded, classID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "delete class success"})
}
