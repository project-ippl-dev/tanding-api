package test

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	usecase Usecase
}

func RegisterHandler(usecase Usecase, e *echo.Echo) {
	handler := handler{usecase: usecase}
	e.GET("/test", handler.test)
}

func (h handler) test(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "success testing"})
}
