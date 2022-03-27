package main

import (
	"io"
	"net/http"
	"video_server/schedular/dbops"

	"github.com/julienschmidt/httprouter"
)

func vidDelRecHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")

	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	io.WriteString(w, "Delete success")
}
