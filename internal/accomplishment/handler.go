package accomplishment

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, profile *echo.Group) {
	handler := handler{usecase: usecase}

	profile.POST("/:uuid/accomplishment", handler.store)
	profile.GET("/:uuid/accomplishment", handler.fetchAll)
	profile.PUT("/:uuid/accomplishment/:accomplishment", handler.update)
	profile.DELETE("/:uuid/accomplishment/:accomplishment", handler.delete)
}

func (h handler) store(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	userID := c.Param("uuid")
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	ctx := c.Request().Context()
	if err := h.usecase.store(ctx, req, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "store accomplishment to specific user success"})
}

func (h handler) fetchAll(c echo.Context) error {
	decoded := tools.JWTDecode(c)
	userID := c.Param("uuid")
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchAll(ctx, page, pageSize, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all accomplishment success", pagination))
}

func (h handler) update(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
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
	if err := h.usecase.update(ctx, req, userID, accomplishmentID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update accomplishment for specific user success"})
}

func (h handler) delete(c echo.Context) error {
	decoded := tools.JWTDecode(c)
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
	if err := h.usecase.delete(ctx, userID, accomplishmentID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete accomplishment for specific user success"})
}
