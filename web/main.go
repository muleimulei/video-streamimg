package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)
	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)

	router.POST("/api", apiHandler)

	router.POST("/upload/:vid-id", proxyUploadHandler)

	router.ServeFiles("/static/*filepath", http.Dir("../templates"))

	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8080", r)
}
