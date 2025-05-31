package club

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

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
	clubHandler := NewHandler(usecase, jwtClient)

	e.POST("/club", clubHandler.Store, m.JWTMiddleware())
	e.PUT("/club/:club", clubHandler.Update, m.JWTMiddleware())
	e.DELETE("/club/:club", clubHandler.Delete, m.JWTMiddleware())
	e.POST("/club/:club/invite", clubHandler.Invite, m.JWTMiddleware())
	e.POST("/club/:club/join", clubHandler.Join, m.JWTMiddleware())
	e.GET("/club/:club/join/approval", clubHandler.FetchJoinApproval, m.JWTMiddleware())
	e.GET("/club/invite/approval", clubHandler.FetchInviteApproval, m.JWTMiddleware())
	e.PATCH("/club/:club/join/approval/:approval", clubHandler.UpdateJoinApproval, m.JWTMiddleware())
	e.PATCH("/club/invite/approval/:approval", clubHandler.UpdateInviteApproval, m.JWTMiddleware())
}

func (h Handler) Store(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation ", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	clubID, err := h.usecase.Store(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, tools.ResponseData{Message: "store club success", Data: clubID})
}

func (h Handler) Update(c echo.Context) error {
	clubParam := c.Param("club")
	clubID, err := uuid.Parse(clubParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	var req Request
	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err = req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	if err = h.usecase.Update(ctx, req, decoded.ID, clubID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update club success"})
}

func (h Handler) Delete(c echo.Context) error {
	clubParam := c.Param("club")
	clubID, err := uuid.Parse(clubParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	if err = h.usecase.Delete(ctx, decoded.ID, clubID); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "delete club success"})
}

func (h Handler) Invite(c echo.Context) error {
	var req ParticipantReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	if err := h.usecase.Invite(ctx, req, decoded.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "success invite participant to club"})
}

func (h Handler) Join(c echo.Context) error {
	var req JoinParam
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Join(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "apply join club success"})
}

func (h Handler) FetchJoinApproval(c echo.Context) error {
	clubIDParam := c.Param("club")
	clubID, err := uuid.Parse(clubIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	limit := tools.PaginationLimit(c)
	ctx := c.Request().Context()
	IDParam := c.QueryParam("id")
	ID, _ := strconv.ParseInt(IDParam, 10, 64)
	data, err := h.usecase.FetchJoinApproval(ctx, limit, clubID, ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch join approval success", Data: data})
}

func (h Handler) FetchInviteApproval(c echo.Context) error {
	limit := tools.PaginationLimit(c)
	IDParam := c.QueryParam("id")
	ID, _ := strconv.ParseInt(IDParam, 10, 64)
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	data, err := h.usecase.FetchInviteApproval(ctx, limit, decoded.ID, ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch invite approval success", Data: data})
}

func (h Handler) UpdateJoinApproval(c echo.Context) error {
	var req UpdateJoinApprovalArgs
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.UpdateJoinApproval(ctx, decoded.ID, req)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update join approval success"})
}

func (h Handler) UpdateInviteApproval(c echo.Context) error {
	var req UpdateInviteApprovalArgs
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	if err := h.usecase.UpdateInviteApproval(ctx, decoded.ID, req); err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.Response{Message: "update invite approval success"})
}
