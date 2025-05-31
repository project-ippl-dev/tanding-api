package file

//go:generate mockgen -source=./usecase.go -destination=../../mocks/file/usecase_mock.go

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/project-ippl-dev/tanding-api/config"
)

type Usecase interface {
	Upload(ctx context.Context, req Params) (path string, err error)
	Base64Upload(req Base64Params) (string, error)
	Delete(dir, path string) (statusCode int, err error)
	Base64ToFile(data string) ([]byte, error)
	ReadMIMEType(src multipart.File, fileInformation *FileInformation) error
}

type usecase struct {
	s3Client config.S3Client
}

func NewUsecase(s3Client config.S3Client) Usecase {
	return &usecase{s3Client: s3Client}
}

func (u usecase) Upload(ctx context.Context, req Params) (path string, err error) {
	src, err := req.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileInfo := req.FileInformation
	filePath := fmt.Sprintf("%s/%d_%s", req.FileInformation.Dir, time.Now().Unix(), fileInfo.FileName)
	svc := s3manager.NewUploader(u.s3Client.Sess)
	object := s3manager.UploadInput{
		Bucket:      &u.s3Client.S3Conf.Bucket,
		ACL:         aws.String("public-read"),
		Key:         aws.String(filePath),
		Body:        src,
		ContentType: aws.String(fileInfo.MIMEType),
	}

	if _, err = svc.Upload(&object); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", u.s3Client.S3Conf.PublicURL, u.s3Client.S3Conf.Bucket, filePath), nil
}

func (u usecase) Base64Upload(req Base64Params) (string, error) {
	fileInfo := req.FileInformation
	filePath := fmt.Sprintf("%s/%s", req.FileInformation.Dir, fileInfo.FileName)
	svc := s3manager.NewUploader(u.s3Client.Sess)
	object := s3manager.UploadInput{
		Bucket:      &u.s3Client.S3Conf.Bucket,
		ACL:         aws.String("public-read"),
		Key:         aws.String(filePath),
		Body:        req.File,
		ContentType: aws.String(fileInfo.MIMEType),
	}

	if _, err := svc.Upload(&object); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", u.s3Client.S3Conf.PublicURL, u.s3Client.S3Conf.Bucket, filePath), nil
}

func (u usecase) Delete(dir, path string) (statusCode int, err error) {
	svc := s3.New(u.s3Client.Sess)
	if _, err = svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(u.s3Client.S3Conf.Bucket),
		Key:    aws.String(dir + "/" + path),
	}); err != nil {
		return http.StatusNotFound, fmt.Errorf("file not found")
	}
	object := s3.DeleteObjectInput{
		Bucket: aws.String(u.s3Client.S3Conf.Bucket),
		Key:    aws.String(dir + "/" + path),
	}
	if _, err = svc.DeleteObject(&object); err != nil {
		return http.StatusRequestTimeout, err
	}
	return http.StatusOK, nil
}

func (u usecase) Base64ToFile(data string) ([]byte, error) {
	i := strings.Index(data, ",")
	raw := data[i+1:]
	return base64.StdEncoding.DecodeString(raw)
}

func (u usecase) ReadMIMEType(src multipart.File, fileInformation *FileInformation) error {
	file := make([]byte, 512)
	if _, err := src.Read(file); err != nil {
		return err
	}
	fileInformation.MIMEType = http.DetectContentType(file)
	fileInformation.Ext = strings.Split(fileInformation.MIMEType, "/")[1]
	return nil
}
