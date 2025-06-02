package testutils

import (
	"bytes"
	"mime/multipart"
	"testing"
)

func GenerateMockFile(t *testing.T, fileName string) (bodyBuf bytes.Buffer, mw *multipart.Writer) {
	mw = multipart.NewWriter(&bodyBuf)

	part, err := mw.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("CreateFormFile failed: %v", err)
	}
	_, _ = part.Write([]byte("hello world")) // any dummy content

	if err = mw.WriteField("dir", "photo"); err != nil {
		t.Fatalf("WriteField failed: %v", err)
	}

	if err = mw.Close(); err != nil {
		t.Fatalf("Close multipart writer: %v", err)
	}

	return
}

func GenerateDummyFile(sizeInMB int) []byte {
	data := make([]byte, sizeInMB*1024*1024)

	for i := range data {
		data[i] = 'A'
	}

	return data
}
