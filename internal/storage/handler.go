package storage

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/project-ippl-dev/tanding-api/internal/middleware"
	"github.com/project-ippl-dev/tanding-api/internal/tools"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) Handler {
	return Handler{
		usecase: usecase,
	}
}

func RegisterHandler(usecase Usecase, m middleware.Params, e *echo.Echo) {
	storageHandler := NewHandler(usecase)

	e.POST("/storage/upload", storageHandler.Upload, m.JWTMiddleware())
	e.POST("/storage/upload/base64", storageHandler.Base64, m.JWTMiddleware())
	e.DELETE("/storage/delete/:dir/:file", storageHandler.Delete)
}

func (h Handler) Upload(c echo.Context) error {
	var req Params
	dir := c.FormValue("dir")
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	req.FileInformation.Dir = dir
	fileName := strings.Replace(file.Filename, " ", "_", -1)
	req.FileInformation.FileName = fileName
	req.FileInformation.Size = file.Size
	req.File = file
	src, err := req.File.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	if err = h.usecase.ReadMIMEType(src, &req.FileInformation); err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	src.Close()
	if err = req.FileInformation.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	ctx := c.Request().Context()
	path, err := h.usecase.Upload(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "upload success", Data: path})
}

func (h Handler) Base64(c echo.Context) error {
	var req Base64Params
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}

	unbased, err := h.usecase.Base64ToFile(req.Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, tools.Response{Message: err.Error()})
	}
	file := bytes.NewReader(unbased)
	mimeType := http.DetectContentType(unbased)
	req.FileInformation.MIMEType = mimeType
	req.FileInformation.FileName = fmt.Sprintf("%d", time.Now().Unix())
	req.FileInformation.Size = file.Size()
	req.FileInformation.Ext = strings.Split(mimeType, "/")[1]
	req.FileInformation.Dir = req.Dir
	req.File = file

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, tools.ResponseValidation{Message: "error validation", Errors: err.Error()})
	}
	path, err := h.usecase.Base64Upload(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, tools.Response{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, tools.ResponseData{Message: "upload success", Data: path})
}

func (h Handler) Delete(c echo.Context) error {
	dir := c.Param("dir")
	path := c.Param("file")
	statusCode, err := h.usecase.Delete(dir, path)
	if err != nil {
		return c.JSON(statusCode, tools.Response{Message: err.Error()})
	}
	return c.JSON(statusCode, tools.Response{Message: "delete file success"})
}
