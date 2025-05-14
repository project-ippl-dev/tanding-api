package mail

import (
	"bytes"
	"html/template"
	"log"

	"github.com/project-ippl-dev/tanding-api/config"
	"gopkg.in/gomail.v2"
)

// Request struct
type Request struct {
	To      string
	Subject string
	Body    bytes.Buffer
}

type MailClient interface {
	SendMail(req Request)
}

type mailClient struct {
	dialer         *gomail.Dialer
	SenderMailName string
}

// MailConnection function is config to setup smtp connection based on environment value
func NewMailClient(mailConf config.MailConfig) MailClient {
	dialer := gomail.NewDialer(mailConf.Host, mailConf.Port, mailConf.Username, mailConf.Password)

	return &mailClient{
		dialer:         dialer,
		SenderMailName: mailConf.Sender,
	}
}

// SendMail represent send mail to specify mail based on parameters
func (ths *mailClient) SendMail(req Request) {
	mailer := gomail.NewMessage()
	mailer.SetHeaders(map[string][]string{
		"From":    {ths.GetSenderMailName()},
		"To":      {req.To},
		"Subject": {req.Subject},
	})
	mailer.SetBody("text/html", req.Body.String())
	if err := ths.dialer.DialAndSend(mailer); err != nil {
		log.Println("error in send mail : " + err.Error())
	} else {
		log.Println("success send mail.")
	}
}

func (ths *mailClient) GetSenderMailName() string {
	return ths.SenderMailName
}

// TemplateToBuffer represent conversion from html page to buffer file
func TemplateToBuffer(templatePath string, data interface{}) (bytes.Buffer, error) {
	var body bytes.Buffer
	templateResult, err := template.ParseFiles(templatePath)
	if err != nil {
		return bytes.Buffer{}, err
	}

	if err := templateResult.Execute(&body, data); err != nil {
		return bytes.Buffer{}, err
	}

	return body, nil
}
