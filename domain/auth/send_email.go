package auth

import (
	"fmt"
	"github.com/lripardo/lrw/domain/api"
)

const (
	LinkMask = "link"
	NameMask = "name"
)

type SendEmailOptions struct {
	Sender   api.EmailSender
	User     *User
	Key      []byte
	Param    string
	Audience string
	Issuer   string
	Route    string
	From     string
	Template string
	Expires  int
}

func sendEmail(options *SendEmailOptions) {
	claims := CreateClaims(options.User.Email, options.Audience, options.Issuer, options.Expires)
	token, err := Sign(claims, options.Key)
	if err != nil {
		api.E("cannot sign claims", err)
		return
	}
	fullLink := fmt.Sprintf("%s?%s=%s", options.Route, options.Param, token)
	templateData := map[string]string{
		NameMask: options.User.FirstName,
		LinkMask: fullLink,
	}
	message := &api.Message{
		From:     options.From,
		To:       options.User.Email,
		Template: options.Template,
		Data:     templateData,
	}
	if err := options.Sender.Send(message); err != nil {
		api.E("cannot send email", err)
	}
	api.D("send email success", options.User.Email)
}
