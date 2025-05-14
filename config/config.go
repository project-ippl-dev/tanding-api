package config

import (
	"github.com/kelseyhightower/envconfig"
	"log"

	"github.com/joho/godotenv"
)

// Config used in this application
type Config struct {
	ServerConfig
	Database      DatabaseConfig `envconfig:"database"`
	JWT           JWTConfig      `envconfig:"jwt"`
	S3            S3Config       `envconfig:"s3"`
	SMTP          MailConfig     `envconfig:"smtp"`
	Redis         RedisConfig    `envconfig:"redis"`
	StorageConfig StorageConfig  `envconfig:"storage"`
}

// DatabaseConfig is specific configuration used for connect to database
type DatabaseConfig struct {
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
	Name     string `envconfig:"name"`
	Host     string `envconfig:"host"`
	Port     string `envconfig:"port"`
}

// ServerConfig is specific configuration used to running application in server
type ServerConfig struct {
	Port int    `envconfig:"port" default:"9001"`
	Host string `envconfig:"host" default:"localhost"`
	FE   string `envconfig:"fe_base_url" default:"http://localhost:3000"`
}

// S3Config is specific configuration used for connect to AWS S3 Bucket
type S3Config struct {
	Endpoint  string `envconfig:"endpoint"`
	PublicURL string
	AccessKey string `envconfig:"access_key"`
	SecretKey string `envconfig:"secret_key"`
	Region    string `envconfig:"region"`
	Bucket    string `envconfig:"bucket"`
}

// S3Config is specific configuration used for hashed jwt token
type JWTConfig struct {
	SecretKey string `envconfig:"secret_key"`
}

// MailConfig is specific configuration used for mailing
type MailConfig struct {
	Host     string `envconfig:"host"`
	Port     int    `envconfig:"port"`
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
	Sender   string `envconfig:"sender"`
}

// RedisConfig is specific configuration used for database cache
type RedisConfig struct {
	Username string `envconfig:"username"`
	Password string `envconfig:"password"`
	Host     string `envconfig:"host"`
	Port     string `envconfig:"port"`
	DB       string
}

type StorageConfig struct {
	BucketID string `envconfig:"bucket_id"`
}

const (
	OAuthGoogle   string = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	SIAAuthKey    string = "sipmarket"
	OAuthFacebook string = "https://graph.facebook.com/v10.0/me?fields=id%2Cname%2Cemail%2Cpicture.height(500)&access_token="
)

// Configuration is function used for load value in env and parsed it to Config struct
func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	var conf Config
	err = envconfig.Process("", &conf)
	if err != nil {
		log.Fatalf("fail to proceed the config: %v", err)
	}
	return conf
}
