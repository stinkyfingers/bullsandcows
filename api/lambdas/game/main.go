package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	amazon "github.com/stinkyfingers/mastermind/aws"
	"github.com/stinkyfingers/mastermind/game"
)

type Event struct {
	Body           string                `json:"body"`
	RequestContext amazon.RequestContext `json:"requestContext"`
}

type Data struct {
	GameID string `json:"gameId"`
}

type Response struct {
	IsBase64 bool              `json:"isBase64Encoded"`
	Status   int               `json:"statusCode"`
	Headers  map[string]string `json:"headers"`
	Body     string            `json:"body"`
}

type ResponseError struct {
	Error string `json:"error"`
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

	var g game.Game
	if d.GameID != "" {
		// get game & resume
		g, err = game.Get(sess, d.GameID)
	} else {
		// create game

		gPtr := game.NewGame(sess, event.RequestContext.ConnectionID)
		g = *gPtr
		err = g.Put()
	}
	if err != nil {
		postError(sess, event.RequestContext.ConnectionID, err)
		return r, nil
	}
	j, err := json.Marshal(g)
	if err != nil {
		return r, err
	}
	r.Body = string(j)
	r.Status = 200
	err = amazon.PostToConnection(sess, event.RequestContext.ConnectionID, j)
	return r, err
}

func postError(sess *session.Session, connectionID string, err error) error {
	j, err := json.Marshal(ResponseError{
		Error: err.Error(),
	})
	if err != nil {
		return err
	}
	return amazon.PostToConnection(sess, connectionID, j)
}
