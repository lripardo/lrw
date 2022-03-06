package email

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/lripardo/lrw/domain/api"
)

var (
	SESName = api.NewKey("AWS_SES_NAME", "required", "Luiz Ricardo Ripardo")
)

type SESEmailSender struct {
	name    string
	service *ses.SES
}

func (s *SESEmailSender) Send(message *api.Message) error {
	from := fmt.Sprintf("%s <%s>", s.name, message.From)
	templateData, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}
	data := string(templateData)

	input := &ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(message.To),
			},
		},
		Source:       aws.String(from),
		Template:     aws.String(message.Template),
		TemplateData: aws.String(data),
	}
	if _, err := s.service.SendTemplatedEmail(input); err != nil {
		return err
	}
	return nil
}

func NewSESEmailSender(configuration api.Configuration, sess *session.Session) api.EmailSender {
	api.D("getting aws ses email sender implementation")
	svc := ses.New(sess)
	name := configuration.String(SESName)
	return &SESEmailSender{
		name:    name,
		service: svc,
	}
}
