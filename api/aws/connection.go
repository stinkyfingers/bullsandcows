package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type RequestContext struct {
	EventType    string `json:"eventType"`
	RequestID    string `json:"requestId"`
	ConnectionID string `json:"connectionId"`
}

type Connection struct {
	ConnectionID string           `json:"connectionId"`
	Channel      string           `json:"channel"`
	Sess         *session.Session `json:"-"`
}

const (
	endpoint = "https://eqsjdaluec.execute-api.us-west-1.amazonaws.com/dev" // connection url
)

func PostToConnection(sess *session.Session, connectionID string, data []byte) error {
	gatewaySess := apigatewaymanagementapi.New(sess, &aws.Config{
		Endpoint: aws.String(endpoint),
	})
	req, _ := gatewaySess.PostToConnectionRequest(&apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         data,
	})
	err := req.Send()
	return err
}
