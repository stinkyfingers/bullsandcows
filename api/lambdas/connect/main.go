package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	amazon "github.com/stinkyfingers/mastermind/aws"
)

type Event struct {
	Body           string                `json:"body"`
	RequestContext amazon.RequestContext `json:"requestContext"`
	Records        interface{}           `json:"records"`
}

type Output struct {
	Message string `json:"message"`
}

type Response struct {
	IsBase64 bool              `json:"isBase64Encoded"`
	Status   int               `json:"statusCode"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
}

type Data struct {
	Action  string `json:"action"`
	Channel string `json:"channel"`
}

const (
	region  = "us-west-1"
	profile = "jds"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event Event) (Response, error) {
	var r Response

	sess, err := amazon.NewAwsSession(region, profile)
	if err != nil {
		return r, err
	}

	switch event.RequestContext.EventType {
	case "CONNECT":
		break

	case "DISCONNECT":
		err = amazon.PostToConnection(sess, event.RequestContext.ConnectionID, []byte("disconnecting: "+event.RequestContext.ConnectionID))
		if err != nil {
			return r, err
		}

	default:
		err = amazon.PostToConnection(sess, event.RequestContext.ConnectionID, []byte("I don't know what to do. Specify an {\"option\":}"))
		if err != nil {
			return r, err
		}
	}

	r.Status = 200
	r.Headers = map[string]string{"testheader": "I'm a header"}
	return r, nil
}
