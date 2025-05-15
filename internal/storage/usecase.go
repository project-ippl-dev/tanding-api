package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/project-ippl-dev/tanding-api/config"
)

const (
	PublicURL = "https://storage.googleapis.com"
)

type Usecase struct {
	storageClient config.StorageClient
}

func NewUsecase(storageClient config.StorageClient) Usecase {
	return Usecase{storageClient: storageClient}
}

func (u Usecase) upload(ctx context.Context, req params) (path string, err error) {
	src, err := req.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileInfo := req.FileInformation
	filePath := fmt.Sprintf("%s/%d_%s", req.FileInformation.Dir, time.Now().Unix(), fileInfo.FileName)
	bucketName := u.storageClient.Config.BucketID

	object := u.storageClient.Client.Bucket(bucketName).Object(filePath)
	wc := object.NewWriter(ctx)
	if _, err = io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	if err = wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", PublicURL, bucketName, filePath), nil
}

func (u Usecase) base64Upload(req base64Params) (string, error) {
	ctx := context.Background()
	fileInfo := req.FileInformation
	filePath := fmt.Sprintf("%s/%s", req.FileInformation.Dir, fileInfo.FileName)
	bucketName := u.storageClient.Config.BucketID

	object := u.storageClient.Client.Bucket(bucketName).Object(filePath)
	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, req.File); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", PublicURL, bucketName, filePath), nil
}

func (u Usecase) Delete(dir, path string) (statusCode int, err error) {
	ctx := context.Background()
	filePath := fmt.Sprintf("%s/%s", dir, path)
	bucketName := u.storageClient.Config.BucketID

	object := u.storageClient.Client.Bucket(bucketName).Object(filePath)
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("object.Attrs: %w", err)
	}
	object = object.If(storage.Conditions{GenerationMatch: attrs.Generation})
	if err = object.Delete(ctx); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Object(%s).Delete: %s", filePath, err.Error())
	}

	return http.StatusOK, nil
}

func (u Usecase) base64ToFile(data string) ([]byte, error) {
	i := strings.Index(data, ",")
	raw := data[i+1:]
	return base64.StdEncoding.DecodeString(raw)
}

func (u Usecase) readMIMEType(src multipart.File, fileInformation *fileInformation) error {
	file := make([]byte, 512)
	if _, err := src.Read(file); err != nil {
		return err
	}
	fileInformation.MIMEType = http.DetectContentType(file)
	fileInformation.Ext = strings.Split(fileInformation.MIMEType, "/")[1]
	return nil
}
