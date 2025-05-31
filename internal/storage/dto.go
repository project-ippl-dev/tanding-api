package storage

import (
	"bytes"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Params struct {
	FileInformation FileInformation
	File            *multipart.FileHeader
}

type Base64Params struct {
	Data            string `json:"data"`
	Dir             string `json:"dir"`
	File            *bytes.Reader
	FileInformation FileInformation
}

func (b Base64Params) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Data, validation.Required),
		validation.Field(&b.Dir, validation.Required),
	)
}

type FileInformation struct {
	//Validation validationRule
	FileName string
	Size     int64
	MIMEType string
	Ext      string
	Dir      string `json:"dir"`
}

func (f FileInformation) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Size, validation.Max(5*1024*1024)),
		validation.Field(&f.Dir, validation.Required),
	)
}

//func checkImageMIMEType(value interface{}) error {
//	val := value.(string)
//	switch val {
//	case "image/jpeg", "image/png":
//		return nil
//	default:
//		return fmt.Errorf("only accept image file")
//	}
//}
