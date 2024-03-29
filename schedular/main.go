package main

import (
	"net/http"
	"video_server/schedular/taskrunner"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)

	return router
}

func main() {
	go taskrunner.Start()
	r := RegisterHandlers()

	http.ListenAndServe(":9001", r)
}
