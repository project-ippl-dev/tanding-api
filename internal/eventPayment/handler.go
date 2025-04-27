package eventPayment

import (
	"github.com/dytlan/tanding-api/internal/middleware"
	"github.com/dytlan/tanding-api/internal/tools"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	handler := handler{usecase: usecase}

	//Role Admin Only
	e.PATCH("/event/:event/payment/:payment", handler.update, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/event/payment", handler.fetchAll, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)
	e.GET("/event/payment/:payment/detail", handler.detail, m.JWTMiddleware(), m.Middleware.RoleAdminOnly)

	//Event Payment For Event Privilege Role Admin
	e.GET("/event/:event/payment/summary", handler.summary, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin)
	e.GET("/event/:event/payment", handler.fetchByEventPrivilege, m.JWTMiddleware(), m.Middleware.EventPrivilegeAdmin)

	//Event Payment For Club Owner
	e.POST("/event/:event/payment", handler.store, m.JWTMiddleware(), m.Middleware.EventManipulationRemarkOpen)
	e.GET("/event/:event/payment/club", handler.fetchByUserID, m.JWTMiddleware())
	e.GET("/event/cart", handler.cart, m.JWTMiddleware())
	e.GET("/event/:event/cart/detail", handler.cartDetail, m.JWTMiddleware())

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
	statusCode, err := h.usecase.store(ctx, req, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "store event payment success"})
}

func (h handler) fetchAll(c echo.Context) error {
	var arg fetchAllParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()

	pagination, err := h.usecase.fetchAll(ctx, page, pageSize, arg)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch all payment receipt success", pagination))
}

func (h handler) update(c echo.Context) error {
	var arg updateParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	if err := arg.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	statusCode, err := h.usecase.update(ctx, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "update payment for admin success"})
}

func (h handler) fetchByUserID(c echo.Context) error {
	var arg fetchByUserIDParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	statusCode, pagination, err := h.usecase.fetchByUserID(ctx, page, pageSize, arg, decoded.ID)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.PaginationGetResponse("fetch by user id success", pagination))
}

func (h handler) fetchByEventPrivilege(c echo.Context) error {
	var arg fetchByEventPrivilegeParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	if arg.EventID == "" {
		return c.JSON(http.StatusUnprocessableEntity, tools.Response{Message: "event param is required"})
	}
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.fetchAll(ctx, page, pageSize, fetchAllParams{
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

func (h handler) cart(c echo.Context) error {
	decoded := tools.JWTDecode(c)
	page, pageSize := tools.PaginationPageAndPageSize(c)
	ctx := c.Request().Context()
	pagination, err := h.usecase.cart(ctx, page, pageSize, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.PaginationGetResponse("fetch cart by specific user success", pagination))
}

func (h handler) cartDetail(c echo.Context) error {
	eventIDParam := c.Param("event")
	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	decoded := tools.JWTDecode(c)
	ctx := c.Request().Context()
	result, err := h.usecase.cartDetail(ctx, eventID, decoded.ID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch cart detail success", Data: result})
}

func (h handler) detail(c echo.Context) error {
	paymentIDParam := c.Param("payment")
	paymentID, err := uuid.Parse(paymentIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	payment, err := h.usecase.detail(ctx, paymentID)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch detail payment success", Data: payment})
}

func (h handler) summary(c echo.Context) error {
	var arg summaryParams
	if err := c.Bind(&arg); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	ctx := c.Request().Context()
	response, err := h.usecase.summary(ctx, arg)
	if err != nil {
		return c.JSON(http.StatusNotFound, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "fetch payment summary success", Data: response})
}
