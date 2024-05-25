package ses_mailer

import (
	"AlbumCodeMailer/mailer"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type SesMailer interface {
	UploadTemplate(*mailer.Template) error
	ListTemplates() ([]types.TemplateMetadata, error)
	SendBulkTemplatedEmail(templateName, fromAddress, replyTo string, recipients []mailer.Recipient) error
}

type SesMailerImpl struct {
	client *ses.Client
}

type ReplacementTemplateData struct {
	Code string `json:"code" validate:"required"`
}

func (s *SesMailerImpl) SendBulkTemplatedEmail(templateName, fromAddress, replyTo string, recipients []mailer.Recipient) error {
	var destinations []types.BulkEmailDestination
	for _, recipient := range recipients {
		rtd := ReplacementTemplateData{
			Code: recipient.AlbumCode,
		}
		s, err := json.Marshal(rtd)
		if err != nil {
			return err
		}
		destinations = append(destinations, types.BulkEmailDestination{
			Destination: &types.Destination{
				ToAddresses: []string{recipient.Email},
			},
			ReplacementTemplateData: aws.String(string(s)),
		})
	}
	sbtei := ses.SendBulkTemplatedEmailInput{
		Destinations:        destinations,
		Source:              aws.String(fromAddress),
		Template:            aws.String(templateName),
		ReturnPath:          aws.String(replyTo),
		DefaultTemplateData: aws.String("{\"code\": \"ERROR\"}"),
	}
	ret, err := s.client.SendBulkTemplatedEmail(context.TODO(), &sbtei)
	if err != nil {
		return err
	}
	for _, d := range ret.Status {
		if d.Error != nil {
			return errors.New(*d.Error)
		}
	}
	return nil
}

func (s *SesMailerImpl) ListTemplates() ([]types.TemplateMetadata, error) {
	ret, err := s.client.ListTemplates(context.TODO(), &ses.ListTemplatesInput{})
	if err != nil {
		return nil, err
	}
	return ret.TemplatesMetadata, nil
}

func (s *SesMailerImpl) UploadTemplate(template *mailer.Template) error {
	var h, pt *string
	if template.HtmlPart != "" {
		h = &template.HtmlPart
	}
	if template.TextPart != "" {
		pt = &template.TextPart
	}
	t := &types.Template{
		TemplateName: &template.Name,
		HtmlPart:     h,
		TextPart:     pt,
		SubjectPart:  &template.Subject,
	}
	_, err := s.client.CreateTemplate(context.TODO(), &ses.CreateTemplateInput{
		Template: t,
	})
	return err
}

func NewSesMailer(region string) (SesMailer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &SesMailerImpl{
		client: ses.NewFromConfig(cfg),
	}, nil
}
