package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config used in this application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	S3       S3Config
	SMTP     MailConfig
	Redis    RedisConfig
}

// DatabaseConfig is specific configuration used for connect to database
type DatabaseConfig struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
}

// ServerConfig is specific configuration used to running application in server
// ServerConfig is specific configuration used to running application in server
type ServerConfig struct {
	Port string
	Host string
	FE   string `mapstructure:"fe"`
}

// S3Config is specific configuration used for connect to AWS S3 Bucket

// S3Config is specific configuration used for connect to AWS S3 Bucket
type S3Config struct {
	Endpoint  string
	PublicURL string
	AccessKey string `mapstructure:"access-key"`
	SecretKey string `mapstructure:"secret-key"`
	Region    string
	Bucket    string
}

// S3Config is specific configuration used for hashed jwt token
type JWTConfig struct {
	SecretKey string
}

// MailConfig is specific configuration used for mailing
type MailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Sender   string
}

// RedisConfig is specific configuration used for database cache
type RedisConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DB       string
}

const (
	OAuthGoogle   string = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	SIAAuthKey    string = "sipmarket"
	OAuthFacebook string = "https://graph.facebook.com/v10.0/me?fields=id%2Cname%2Cemail%2Cpicture.height(500)&access_token="
)

// Configuration is function used for load value in env and parsed it to Config struct
func Configuration() Config {
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("Failed to unmarshall" + err.Error())
	}
	return conf
}
