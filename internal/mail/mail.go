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

// SendMail represent send mail to specify mail based on parameters
func SendMail(req Request) {
	mailer := gomail.NewMessage()
	mailer.SetHeaders(map[string][]string{
		"From":    {config.GetSenderMailName()},
		"To":      {req.To},
		"Subject": {req.Subject},
	})
	mailer.SetBody("text/html", req.Body.String())
	dialer := config.MailConnection()
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Println("error in send mail : " + err.Error())
	} else {
		log.Println("success send mail.")
	}
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
