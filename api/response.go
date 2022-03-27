package main

import (
	"encoding/json"
	"io"
	"net/http"
	"video_server/api/defs"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrorResponse) {
	w.WriteHeader(errResp.HttpSc)
	data, _ := json.Marshal(errResp.Error)
	io.WriteString(w, string(data))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
