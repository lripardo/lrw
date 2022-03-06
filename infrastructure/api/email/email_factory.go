package email

import (
	"github.com/lripardo/lrw/domain/api"
	"github.com/lripardo/lrw/infrastructure/api/connection"
)

const (
	ConsoleEmailType = "console"
	SESEmailType     = "ses"
)

var (
	Config = api.NewKey("EMAIL_SENDER_TYPE", "required", ConsoleEmailType)
)

func NewEmailService(configuration api.Configuration) (api.EmailSender, error) {
	emailConfig := configuration.String(Config)
	if emailConfig == SESEmailType {
		session, err := connection.CreateAwsSession()
		if err != nil {
			return nil, err
		}
		return NewSESEmailSender(configuration, session), nil
	}
	return NewConsoleEmailSender(), nil
}
