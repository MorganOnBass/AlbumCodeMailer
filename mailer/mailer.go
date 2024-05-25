package mailer

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Recipient struct {
	Name      string
	Email     string
	AlbumCode string
}

type Template struct {
	Name     string `json:"name" validate:"required,string"`
	HtmlPart string `json:"html_part" validate:"string"`
	Subject  string `json:"subject" validate:"required,string"`
	TextPart string `json:"text_part" validate:"string"`
}

func NewTemplateFromTextFiles(htmlpart, textpart io.Reader, name, subject string) (*Template, error) {
	var h, t []byte
	var err error
	if htmlpart != nil {
		h, err = io.ReadAll(htmlpart)
		if err != nil {
			return nil, err
		}
		if len(h) < 1 {
			return nil, errors.New("html part is empty")
		}
	}
	t, err = io.ReadAll(textpart)
	if err != nil {
		return nil, err
	}
	if len(t) < 1 {
		return nil, errors.New("text part is empty")
	}
	return &Template{
		Name:     name,
		Subject:  subject,
		HtmlPart: string(h),
		TextPart: string(t),
	}, nil
}

func NewTemplateFromFile(file string, name string) (*Template, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	t := &Template{}
	err = json.Unmarshal(buf, t)
	return t, err
}
