package email

import "github.com/lripardo/lrw/domain/api"

type ConsoleEmailSender struct{}

func (n *ConsoleEmailSender) Send(message *api.Message) error {
	api.I("Email: ", message.Template, message.From, message.To, message.Data)
	return nil
}

func NewConsoleEmailSender() api.EmailSender {
	api.D("getting console email sender implementation")
	return &ConsoleEmailSender{}
}
