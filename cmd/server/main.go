package main

import (
	"log"

	"github.com/project-ippl-dev/tanding-api/internal/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error load config : " + err.Error())
	}
	if err := cmd.Run(); err != nil {
		log.Fatal("Running Failed : " + err.Error())
	}
}
