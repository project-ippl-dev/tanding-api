package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3Client struct {
	Sess   *session.Session
	S3Conf S3Config
}

// S3Connection is function used to open connection to AWS S3 Bucket
func NewS3Client(s3Conf S3Config) S3Client {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3Conf.AccessKey, s3Conf.SecretKey, ""),
		Region:           &s3Conf.Region,
		Endpoint:         &s3Conf.Endpoint,
		S3ForcePathStyle: aws.Bool(true),
	}))
	return S3Client{
		Sess:   sess,
		S3Conf: s3Conf,
	}
}
