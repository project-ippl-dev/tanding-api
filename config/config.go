package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config used in this application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	S3       S3Config
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
type ServerConfig struct {
	Port    string
	AuthKey string `mapstructure:"auth-key"`
}

// S3Config is specific configuration used for connect to AWS S3 Bucket
type S3Config struct {
	AccessKey string `mapstructure:"access-key"`
	SecretKey string `mapstructure:"secret-key"`
	Region    string
	Bucket    string
}

// S3Config is specific configuration used for hashed jwt token
type JWTConfig struct {
	SecretKey string
}

const (
	OAuthGoogle string = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	SIAAuthKey  string = "sipmarket"
)

// Configuration is function used for load value in env and parsed it to Config struct
func Configuration() Config {
	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal("Failed to unmarshall" + err.Error())
	}
	return conf
}

// AuthKey is function used to get env value of Config.Server.AuthKey
func AuthKey() string {
	return Configuration().Server.AuthKey
}
