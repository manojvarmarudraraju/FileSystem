package handlerFunctions

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"file-stats/internal/web"
	"io/ioutil"
	"log"
	"net/http"
)

type PostHandleVars struct{
	RequestInputChan chan FilesStruct
	Log *log.Logger
}

type File struct {
	FileAttr xml.Name `xml:"file"`
	Filesize  int16 `xml:"filesize",json:"filesize"`
	Extension string `xml:"extension",json:"extension"`
	FilePath  string `xml:"filepath",json:"filepath"`
}

type FilesStruct struct {
	FilesAttr xml.Name`xml:"files,attr"`
	Files []File `xml:"file",json:"files"`
}

// handling the post request and pushes the data into channel to avoid memory access at the same time

func (PHV PostHandleVars) PostFiles(ctx context.Context, w http.ResponseWriter,r *http.Request) error {
	var err  error
	var Files FilesStruct
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	PHV.Log.Println(r.Header["Content-Type"][0])
	if err != nil {
		return err
	}
	if r.Header["Content-Type"][0] == "application/xml"{
		if err = xml.Unmarshal(b,&Files); err != nil{
			PHV.Log.Println("Couldn't decode body XML: ", err)
			return web.RespondError(ctx, w, err, http.StatusBadRequest)
		}
	}else{
		if err = json.Unmarshal(b,&Files); err != nil{
			PHV.Log.Println("Couldn't decode body JSON: ", err)
			return web.RespondError(ctx, w, err, http.StatusBadRequest)
		}
	}
	PHV.RequestInputChan <- Files
	return web.Respond(ctx,w,nil,http.StatusOK)
}
