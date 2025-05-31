package eventPayment

import (
	"net/http"

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
	eventPaymentHandler := NewHandler(usecase, jwtClient)

	//Role Admin Only
	e.PATCH("/event/:event/payment/:payment", eventPaymentHandler.Update, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/event/payment", eventPaymentHandler.FetchAll, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/event/payment/:payment/detail", eventPaymentHandler.Detail, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)

	//Event Payment For Event Privilege Role Admin
	e.GET("/event/:event/payment/summary", eventPaymentHandler.Summary, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin)
	e.GET("/event/:event/payment", eventPaymentHandler.FetchByEventPrivilege, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin)

	//Event Payment For Club Owner
	e.POST("/event/:event/payment", eventPaymentHandler.Store, m.JWTMiddleware(), m.Middleware.EventManipulationRemarkOpen)
	e.GET("/event/:event/payment/club", eventPaymentHandler.FetchByUserID, m.JWTMiddleware())
	e.GET("/event/cart", eventPaymentHandler.Cart, m.JWTMiddleware())
	e.GET("/event/:event/cart/detail", eventPaymentHandler.CartDetail, m.JWTMiddleware())

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
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Store(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "store event payment success"})
}

func (h Handler) FetchAll(c echo.Context) error {
	var arg FetchAllParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()

	pagination, err := h.usecase.FetchAll(ctx, page, pageSize, arg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all payment receipt success", pagination))
}

func (h Handler) Update(c echo.Context) error {
	var arg UpdateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.Update(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update payment for admin success"})
}

func (h Handler) FetchByUserID(c echo.Context) error {
	var arg FetchByUserIDParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	statusCode, pagination, err := h.usecase.FetchByUserID(ctx, page, pageSize, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.PaginationGetResponse("fetch by user id success", pagination))
}

func (h Handler) FetchByEventPrivilege(c echo.Context) error {
	var arg FetchByEventPrivilegeParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if arg.EventID == "" {
		return c.JSON(http.StatusUnprocessableEntity, tools.Response{Message: "event param is required"})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.FetchAll(ctx, page, pageSize, FetchAllParams{
		EventID: arg.EventID,
		Status:  arg.Status,
		Clubs:   arg.Clubs,
		Start:   arg.Start,
		End:     arg.End,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all payment receipt in specific event success", pagination))
}

func (h Handler) Cart(c echo.Context) error {
	decoded := h.jwtClient.Decode(c)
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.Cart(ctx, page, pageSize, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch cart by specific user success", pagination))
}

func (h Handler) CartDetail(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := h.jwtClient.Decode(c)
	ctx := c.Request().Context()
	result, err := h.usecase.CartDetail(ctx, eventID, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch cart detail success", Data: result})
}

func (h Handler) Detail(c echo.Context) error {
	paymentIDParam := c.Param("payment")
	paymentID, err := uuid.Parse(paymentIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	payment, err := h.usecase.Detail(ctx, paymentID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch detail payment success", Data: payment})
}

func (h Handler) Summary(c echo.Context) error {
	var arg SummaryParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	response, err := h.usecase.Summary(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch payment summary success", Data: response})
}
