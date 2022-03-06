package api

type Message struct {
	From     string
	To       string
	Template string
	Data     map[string]string
}

type EmailSender interface {
	Send(message *Message) error
}
