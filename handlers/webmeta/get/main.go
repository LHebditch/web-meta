package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/LHebditch/htmlmeta/handlers/webmeta"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

type GetMetaEvent struct {
	URL string `json:"url"`
}

var log *zap.Logger

func init() {
    logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialise zap logger")
	}
	log = logger
}

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	defer log.Sync()
	eventBody := GetMetaEvent{}
	err := json.Unmarshal([]byte(event.Body), &eventBody)
	if err != nil {
		panic("failed to parse body")
	}
	if eventBody.URL == "" {
		panic("invalid url")
	}
	client := &http.Client{}
	meta, err := webmeta.GetWebMeta(log, eventBody.URL, client)
	if err != nil {
		panic("failed to get web meta from url")
	}
	
	json, _ := json.Marshal(meta)
	
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: string(json),
		IsBase64Encoded: false,
	}, nil
}

func main() {
   lambda.Start(Handler)
}