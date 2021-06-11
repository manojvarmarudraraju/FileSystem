package web

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func RespondWithFile(ctx context.Context, w http.ResponseWriter, r *http.Request, data *bytes.Buffer, statusCode int ,filename string) error  {

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok{
		return errors.New("web values missing from context")
	}
	v.StatusCode=statusCode
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data.Bytes()))

	return nil
}
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int)error  {

	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok{
		return errors.New("web values missing from context")
	}
	v.StatusCode=statusCode
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}
	//once we have the data do marshal into json to send it to ui.
	mData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//set the status type and the content type and text format in proper order, to avoid overwriting.
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if _, err := w.Write(mData); err != nil {
		return err
	}
	return nil
}

//RespondError knows how to handle the errors going out.
func RespondError(ctx context.Context, w http.ResponseWriter, err error,status int)error  {
	resp:= ErrorResponse{Error:err.Error()}
	return Respond(ctx,w,resp,status)
}
