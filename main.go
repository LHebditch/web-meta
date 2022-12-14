package main

import (
	"net/http"

	"github.com/LHebditch/htmlmeta/handlers/webmeta"
	"go.uber.org/zap"
)

func main(){
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialise zap logger")
	}
	client := &http.Client{}
	meta, err := webmeta.GetWebMeta(logger, "https://www.iwantoneofthose.com/gift-home-office/mario-kart-green-shell-light-with-sound/13030414.html", client)
	if err != nil {
		panic(err.Error())
	}

	json, _ := meta.ToString()
	logger.Info(json)
}