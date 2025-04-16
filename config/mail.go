package config

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

//GetSenderMailName function to get environment of smtp sender name
func GetSenderMailName() string {
	return Configuration().SMTP.Sender
}

//MailConnection function is config to setup smtp connection based on environment value
func MailConnection() *gomail.Dialer {
	conf := Configuration().SMTP

	port, _ := strconv.Atoi(conf.Port)

	dialer := gomail.NewDialer(conf.Host, port, conf.Username, conf.Password)

	return dialer
}
