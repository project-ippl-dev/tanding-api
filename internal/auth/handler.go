package auth

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/config"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, e *echo.Echo) {
	handler := handler{usecase: usecase}

	jwtMiddleware := tools.JWTMiddleware()

	auth := e.Group("/auth")
	auth.POST("/register", handler.register)
	auth.GET("/verification/:type/:token", handler.verify)
	auth.POST("/login", handler.login)
	auth.POST("/callback/:type/:token", handler.callback)
	auth.GET("/check/:type/:request", handler.check)
	auth.GET("/forgot-password/:username", handler.forgot)
	auth.POST("/reset-password/:token", handler.reset)
	auth.GET("/resend-token/:type/:username", handler.resend)
	auth.POST("/bind/:type", handler.binding, jwtMiddleware)
}

func (h handler) register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	s := config.Configuration().Server
	statusCode, err := h.usecase.register(ctx, req, s.Host, s.FE)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "register success, check email for verify."})
}

func (h handler) verify(c echo.Context) error {
	token := c.Param("token")
	kind := c.Param("type")
	decoded, err := tools.JWTTokenParse(token)
	if err != nil {
		return c.JSON(http.StatusForbidden, tools.Response{Message: "Token is not valid or expired, error : " + err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, email, err := h.usecase.verify(ctx, kind, decoded, token)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: kind + " account success", Data: email})
}

func (h handler) login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in binding request : " + err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error in validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, response, err := h.usecase.login(ctx, req)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: "login success", Data: response})
}

func (h handler) callback(c echo.Context) error {
	kind := c.Param("type")
	accessToken := c.Param("token")
	ctx := c.Request().Context()
	statusCode, response, err := h.usecase.callback(ctx, kind, accessToken)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: "login via " + kind + " success", Data: response})
}

func (h handler) check(c echo.Context) error {
	kind := c.Param("type")
	request := c.Param("request")
	ctx := c.Request().Context()
	data, err := h.usecase.check(ctx, kind, request)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch account success", Data: data})
}

func (h handler) forgot(c echo.Context) error {
	username := c.Param("username")
	ctx := c.Request().Context()
	s := config.Configuration().Server
	if err := validation.Validate(&username, validation.Required, is.Email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	statusCode, err := h.usecase.forgot(ctx, username, s.Host, s.FE)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "send forgot password mail success"})
}

func (h handler) reset(c echo.Context) error {
	token := c.Param("token")
	var req resetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in binding request : " + err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded, err := tools.JWTTokenParse(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in parsing token, token maybe expired or invalid : " + err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.reset(ctx, req, decoded, token)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "reset password success"})
}

func (h handler) resend(c echo.Context) error {
	kind := c.Param("type")
	username := c.Param("username")
	ctx := c.Request().Context()
	if err := validation.Validate(&username, validation.Required, is.Email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	statusCode, err := h.usecase.resend(ctx, username, kind)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "resend token success"})
}

func (h handler) binding(c echo.Context) error {
	kind := c.Param("type")
	decoded := tools.JWTDecode(c)

	var req bindingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in binding request : " + err.Error()})
	}
	if err := req.Validate(kind); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	ctx := c.Request().Context()
	statusCode, err := h.usecase.binding(ctx, bindingParams{
		Kind:    kind,
		Request: req,
		Decoded: decoded,
		Host:    config.Configuration().Server.Host,
	})
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	switch kind {
	case string(db.AccountTypeManual):
		return c.JSON(statusCode, tools.Response{Message: "check email, for binding is just one step ahead"})
	case string(db.AccountTypeFacebook), string(db.AccountTypeGoogle):
		return c.JSON(statusCode, tools.Response{Message: "binding account via " + kind + " success"})
	}
	return c.JSON(http.StatusInternalServerError, tools.Response{Message: "something went wrong"})
}
