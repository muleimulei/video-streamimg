package main

/*
handler -> validation{1.request 2.user} -> business logic ->response
1. data model
2. error handling
*/

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 检查session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)
	router.POST("/user/:name", Login)
	return router
}

func main() {
	r := RegisterHandlers()
	mw := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mw)
}
