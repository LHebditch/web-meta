package main

import (
	"github.com.LHebditch.htmlmeta/handlers/webmeta"
	"go.uber.org/zap"
)

func main(){
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to initialise zap logger")
	}
	meta, err := webmeta.GetWebMeta(logger, "https://ign.com/uk")
	if err != nil {
		panic(err.Error())
	}

	json, _ := meta.ToString()
	logger.Info(json)
}