package auth

import (
	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	"github.com/project-ippl-dev/tanding-api/config"
	"github.com/project-ippl-dev/tanding-api/internal/db"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Handler struct {
	usecase    Usecase
	jwtClient  tools.JWTClient
	serverConf config.ServerConfig
}

func NewHandler(usecase Usecase, jwtClient tools.JWTClient, serverConf config.ServerConfig) Handler {
	return Handler{
		usecase:    usecase,
		jwtClient:  jwtClient,
		serverConf: serverConf,
	}
}

func RegisterHandler(usecase Usecase, jwtClient tools.JWTClient, m middleware.Params, serverConf config.ServerConfig, e *echo.Echo) {
	authHandler := NewHandler(usecase, jwtClient, serverConf)

	auth := e.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.GET("/verification/:type/:token", authHandler.Verify)
	auth.POST("/login", authHandler.Login)
	auth.POST("/callback/:type/:token", authHandler.Callback)
	auth.GET("/check/:type/:request", authHandler.Check)
	auth.GET("/forgot-password/:username", authHandler.Forgot)
	auth.POST("/reset-password/:token", authHandler.Reset)
	auth.GET("/resend-token/:type/:username", authHandler.Resend)
	auth.POST("/bind/:type", authHandler.Binding, m.JWTMiddleware())
}

func (h Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Register(ctx, req, h.serverConf.Host, h.serverConf.FE)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "register success, check email for verify."})
}

func (h Handler) Verify(c echo.Context) error {
	token := c.Param("token")
	kind := c.Param("type")
	decoded, err := h.jwtClient.TokenParse(token)
	if err != nil {
		return c.JSON(http.StatusForbidden, tools.Response{Message: "Token is not valid or expired, error : " + err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, email, err := h.usecase.Verify(ctx, kind, decoded, token)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: kind + " account success", Data: email})
}

func (h Handler) Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in binding request : " + err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error in validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, response, err := h.usecase.Login(ctx, req)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: "login success", Data: response})
}

func (h Handler) Callback(c echo.Context) error {
	kind := c.Param("type")
	accessToken := c.Param("token")
	ctx := c.Request().Context()
	statusCode, response, err := h.usecase.Callback(ctx, kind, accessToken)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.ResponseData{Message: "login via " + kind + " success", Data: response})
}

func (h Handler) Check(c echo.Context) error {
	kind := c.Param("type")
	request := c.Param("request")
	ctx := c.Request().Context()
	data, err := h.usecase.Check(ctx, kind, request)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch account success", Data: data})
}

func (h Handler) Forgot(c echo.Context) error {
	username := c.Param("username")
	ctx := c.Request().Context()
	if err := validation.Validate(&username, validation.Required, is.Email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	statusCode, err := h.usecase.Forgot(ctx, username, h.serverConf.Host, h.serverConf.FE)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "send forgot password mail success"})
}

func (h Handler) Reset(c echo.Context) error {
	token := c.Param("token")
	var req ResetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in binding request : " + err.Error()})
	}
	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded, err := h.jwtClient.TokenParse(token)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in parsing token, token maybe expired or invalid : " + err.Error()})
	}
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Reset(ctx, req, decoded, token)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "reset password success"})
}

func (h Handler) Resend(c echo.Context) error {
	kind := c.Param("type")
	username := c.Param("username")
	ctx := c.Request().Context()
	if err := validation.Validate(&username, validation.Required, is.Email); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	statusCode, err := h.usecase.Resend(ctx, username, kind)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "resend token success"})
}

func (h Handler) Binding(c echo.Context) error {
	kind := c.Param("type")
	decoded := h.jwtClient.Decode(c)

	var req BindingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: "error in binding request : " + err.Error()})
	}
	if err := req.Validate(kind); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}

	ctx := c.Request().Context()
	statusCode, err := h.usecase.Binding(ctx, BindingParams{
		Kind:    kind,
		Request: req,
		Decoded: decoded,
		Host:    h.serverConf.Host,
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
