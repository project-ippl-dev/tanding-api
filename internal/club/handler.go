package club

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	handler := handler{usecase: usecase}

	e.POST("/club", handler.store, m.JWTMiddleware())
	e.PUT("/club/:club", handler.update, m.JWTMiddleware())
	e.DELETE("/club/:club", handler.delete, m.JWTMiddleware())
	e.POST("/club/:club/invite", handler.invite, m.JWTMiddleware())
	e.POST("/club/:club/join", handler.join, m.JWTMiddleware())
	e.GET("/club/:club/join/approval", handler.fetchJoinApproval, m.JWTMiddleware())
	e.GET("/club/invite/approval", handler.fetchInviteApproval, m.JWTMiddleware())
	e.PATCH("/club/:club/join/approval/:approval", handler.updateJoinApproval, m.JWTMiddleware())
	e.PATCH("/club/invite/approval/:approval", handler.updateInviteApproval, m.JWTMiddleware())
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
	ctx := c.Request().Context()
	clubID, err := h.usecase.store(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.ResponseData{Message: "store club success", Data: clubID})
}

func (h handler) update(c echo.Context) error {
	clubParam := c.Param("club")
	clubID, err := uuid.Parse(clubParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	var req request
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err = req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	if err = h.usecase.update(ctx, req, decoded.ID, clubID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update club success"})
}

func (h handler) delete(c echo.Context) error {
	clubParam := c.Param("club")
	clubID, err := uuid.Parse(clubParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	if err = h.usecase.delete(ctx, decoded.ID, clubID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete club success"})
}

func (h handler) invite(c echo.Context) error {
	var req participantReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	if err := h.usecase.invite(ctx, req, decoded.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "success invite participant to club"})
}

func (h handler) join(c echo.Context) error {
	var req joinParam
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.join(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "apply join club success"})
}

func (h handler) fetchJoinApproval(c echo.Context) error {
	clubIDParam := c.Param("club")
	clubID, err := uuid.Parse(clubIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	limit := tools.PaginationLimit(c)
	ctx := c.Request().Context()
	IDParam := c.QueryParam("id")
	ID, _ := strconv.ParseInt(IDParam, 10, 64)
	data, err := h.usecase.fetchJoinApproval(ctx, limit, clubID, ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch join approval success", Data: data})
}

func (h handler) fetchInviteApproval(c echo.Context) error {
	limit := tools.PaginationLimit(c)
	IDParam := c.QueryParam("id")
	ID, _ := strconv.ParseInt(IDParam, 10, 64)
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	data, err := h.usecase.fetchInviteApproval(ctx, limit, decoded.ID, ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch invite approval success", Data: data})
}

func (h handler) updateJoinApproval(c echo.Context) error {
	var req updateJoinApprovalArgs
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.updateJoinApproval(ctx, decoded.ID, req)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update join approval success"})
}

func (h handler) updateInviteApproval(c echo.Context) error {
	var req updateInviteApprovalArgs
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	if err := h.usecase.updateInviteApproval(ctx, decoded.ID, req); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update invite approval success"})
}
