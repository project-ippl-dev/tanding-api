package accomplishment

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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

func RegisterHandler(usecase Usecase, jwtClient tools.JWTClient, profile *echo.Group) {
	accomplishmentHandler := NewHandler(usecase, jwtClient)

	profile.POST("/:uuid/accomplishment", accomplishmentHandler.Store)
	profile.GET("/:uuid/accomplishment", accomplishmentHandler.FetchAll)
	profile.PUT("/:uuid/accomplishment/:accomplishment", accomplishmentHandler.Update)
	profile.DELETE("/:uuid/accomplishment/:accomplishment", accomplishmentHandler.Delete)
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
	userID := c.Param("uuid")
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	ctx := c.Request().Context()
	if err := h.usecase.Store(ctx, req, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "store accomplishment to specific user success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	decoded := h.jwtClient.Decode(c)
	userID := c.Param("uuid")
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all accomplishment success", pagination))
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
	userID := c.Param("uuid")
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	accomplishmentIDParam := c.Param("accomplishment")
	accomplishmentID, err := strconv.ParseInt(accomplishmentIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err = h.usecase.Update(ctx, req, userID, accomplishmentID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update accomplishment for specific user success"})
}

func (h Handler) Delete(c echo.Context) error {
	decoded := h.jwtClient.Decode(c)
	userID := c.Param("uuid")
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	accomplishmentIDParam := c.Param("accomplishment")
	accomplishmentID, err := strconv.ParseInt(accomplishmentIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err = h.usecase.Delete(ctx, userID, accomplishmentID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete accomplishment for specific user success"})
}
