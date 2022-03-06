package connection

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/lripardo/lrw/domain/api"
)

var awsSession *session.Session

func CreateAwsSession() (*session.Session, error) {
	if awsSession == nil {
		s, err := session.NewSession()
		if err != nil {
			return nil, err
		}
		awsSession = s
	}
	api.D("getting new aws session")
	return awsSession, nil
}
