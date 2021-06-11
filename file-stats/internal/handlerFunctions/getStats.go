package handlerFunctions

import (
	"context"
	"file-stats/internal/web"
	"fmt"
	"log"
	"net/http"
)

type Stats struct {
	FilesCount int64 `xml:"filesCount,attr",json:"filesCount"`
	MaxFileSize int64 `xml:"maxfilesize,attr",json:"maxfilesize"`
	AverageFileSize float64 `xml:"avgFilesize,attr",json:"avgFilesize"`
	ListExtensions []string `xml:"listExtensions,attr",json:"listExtensions"`
	FileExtensionsCount map[string]int64 `xml:"fileExtensionsCount,attr",json:"fileExtensionsCount"`
	LatestPath []string `xml:"latestPaths,attr",json:"latestPaths"`
}

type GetStatsVars struct{
	Log *log.Logger
	Stats *Stats
}

// get handler to provide stats

func (GSV GetStatsVars) GetStats(ctx context.Context, w http.ResponseWriter,r *http.Request) error {
	fmt.Println("Came here")
	return web.Respond(ctx,w,GSV.Stats,http.StatusOK)
}
