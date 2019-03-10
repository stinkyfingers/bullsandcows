package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	amazon "github.com/stinkyfingers/mastermind/aws"
	"github.com/stinkyfingers/mastermind/game"
)

type Event struct {
	Body           string                `json:"body"`
	RequestContext amazon.RequestContext `json:"requestContext"`
}

type Data struct {
	GameID   string        `json:"gameId"`
	Sequence game.Sequence `json:"sequence"`
}

type Response struct {
	IsBase64 bool              `json:"isBase64Encoded"`
	Status   int               `json:"statusCode"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
}

const (
	region  = "us-west-1"
	profile = "jds"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event Event) (Response, error) {
	r := Response{Status: 500}
	sess, err := amazon.NewAwsSession(region, profile)
	if err != nil {
		return r, err
	}

	var d Data
	err = json.Unmarshal([]byte(event.Body), &d)
	if err != nil {
		return r, err
	}

	g, err := game.Get(sess, d.GameID)
	if err != nil {
		return r, err
	}

	err = g.PlayRound(d.Sequence)
	if err != nil {
		return r, err
	}

	j, err := json.Marshal(g)
	if err != nil {
		return r, err
	}

	r.Status = 200
	err = amazon.PostToConnection(sess, event.RequestContext.ConnectionID, j)
	return r, err
}
