package main

import (
	"context"
	"net/http"

	"github.com.LHebditch.htmlmeta/handlers/webmeta"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

type GetMetaEvent struct {
	URL string
}

var log *zap.Logger

func init() {
    logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialise zap logger")
	}
	log = logger
}

func Handler(ctx context.Context, event GetMetaEvent) (string, error) {
	defer log.Sync()
	if event.URL == "" {
		panic("invalid url")
	}
	client := &http.Client{}
	meta, err := webmeta.GetWebMeta(log, event.URL, client)
	if err != nil {
		panic("failed to get web meta from url")
	}
	json, err := meta.ToString()
	if err != nil {

	}
	return json, nil
}

func main() {
   lambda.Start(Handler)
}