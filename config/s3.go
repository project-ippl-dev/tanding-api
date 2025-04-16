package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// S3Credential is function to get all value of Config.S3Config
func S3Credential() S3Config {
	return Configuration().S3
}

// S3Connection is function used to open connection to AWS S3 Bucket
func S3Connection() *session.Session {
	s3Config := S3Credential()
	endpoint := s3Config.Endpoint
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3Config.AccessKey, s3Config.SecretKey, ""),
		Region:           &s3Config.Region,
		Endpoint:         &endpoint,
		S3ForcePathStyle: aws.Bool(true),
	}))
	return sess
}
