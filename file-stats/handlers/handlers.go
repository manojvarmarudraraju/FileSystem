package handlers

import (
	"file-stats/internal/handlerFunctions"
	"file-stats/internal/mid"
	"file-stats/internal/web"
	"log"
	"net/http"
)

func API(logger *log.Logger,requestInput chan handlerFunctions.FilesStruct,Stats *handlerFunctions.Stats) http.Handler {
	// web app with logger , error and metrics middleware
	app := web.NewApp(logger,mid.Logger(logger),mid.Errors(logger))
	postVar := handlerFunctions.PostHandleVars{
		RequestInputChan: requestInput,
		Log: logger,
	}
	getVar := handlerFunctions.GetStatsVars{
		Log: logger,
		Stats: Stats,
	}
	app.Handle(http.MethodPost, "/files", postVar.PostFiles)
	app.Handle(http.MethodGet, "/Stats", getVar.GetStats)

	return app
}
