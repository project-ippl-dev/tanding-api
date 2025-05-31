package document

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
	documentHandler := NewHandler(usecase, jwtClient)

	profile.POST("/:uuid/document", documentHandler.Store)
	profile.GET("/:uuid/document", documentHandler.FetchAll)
	profile.PUT("/:uuid/document/:document", documentHandler.Update)
	profile.DELETE("/:uuid/document/:document", documentHandler.Delete)
}

func (h Handler) Store(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	userID := c.Param("uuid")
	decoded := h.jwtClient.Decode(c)
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	ctx := c.Request().Context()
	if err := h.usecase.Store(ctx, req, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.Response{Message: "store document to specific user success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	page, pageSize := tools.PaginationPageAndPageSize(c)
	userID := c.Param("uuid")
	decoded := h.jwtClient.Decode(c)
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch document for specific user success", pagination))
}

func (h Handler) Update(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	userID := c.Param("uuid")
	decoded := h.jwtClient.Decode(c)
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	documentIDParam := c.Param("document")
	documentID, err := strconv.ParseInt(documentIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Update(ctx, req, userID, documentID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update document for specific user success"})
}

func (h Handler) Delete(c echo.Context) error {
	userID := c.Param("uuid")
	decoded := h.jwtClient.Decode(c)
	if decoded.RoleName == "user" {
		userID = decoded.ID
	}
	documentIDParam := c.Param("document")
	documentID, err := strconv.ParseInt(documentIDParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if err := h.usecase.Delete(ctx, userID, documentID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete document for specific user success"})
}
