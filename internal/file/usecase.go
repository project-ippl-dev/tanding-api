package file

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/project-ippl-dev/tanding-api/config"
)

type Usecase struct {
	sess *session.Session
}

func NewUsecase(sess *session.Session) Usecase {
	return Usecase{sess: sess}
}

func (u Usecase) upload(ctx context.Context, req params) (path string, err error) {
	src, err := req.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileInfo := req.FileInformation
	s3Conf := config.S3Credential()
	filePath := fmt.Sprintf("%s/%d_%s", req.FileInformation.Dir, time.Now().Unix(), fileInfo.FileName)
	svc := s3manager.NewUploader(u.sess)
	object := s3manager.UploadInput{
		Bucket:      &s3Conf.Bucket,
		ACL:         aws.String("public-read"),
		Key:         aws.String(filePath),
		Body:        src,
		ContentType: aws.String(fileInfo.MIMEType),
	}

	if _, err := svc.Upload(&object); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s3Conf.PublicURL, s3Conf.Bucket, filePath), nil
}

func (u Usecase) base64Upload(req base64Params) (string, error) {
	fileInfo := req.FileInformation
	s3Conf := config.S3Credential()
	filePath := fmt.Sprintf("%s/%s", req.FileInformation.Dir, fileInfo.FileName)
	svc := s3manager.NewUploader(u.sess)
	object := s3manager.UploadInput{
		Bucket:      &s3Conf.Bucket,
		ACL:         aws.String("public-read"),
		Key:         aws.String(filePath),
		Body:        req.File,
		ContentType: aws.String(fileInfo.MIMEType),
	}

	if _, err := svc.Upload(&object); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", s3Conf.PublicURL, s3Conf.Bucket, filePath), nil
}

func (u Usecase) Delete(dir, path string) (statusCode int, err error) {
	svc := s3.New(u.sess)
	if _, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(config.S3Credential().Bucket),
		Key:    aws.String(dir + "/" + path),
	}); err != nil {
		return http.StatusNotFound, fmt.Errorf("file not found")
	}
	object := s3.DeleteObjectInput{
		Bucket: aws.String(config.S3Credential().Bucket),
		Key:    aws.String(dir + "/" + path),
	}
	if _, err := svc.DeleteObject(&object); err != nil {
		return http.StatusRequestTimeout, err
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
