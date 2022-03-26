package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//创建用户
func CreateUser(w http.ResponseWriter, r *http.Request, args httprouter.Params) {
	io.WriteString(w, args.ByName("name"))
}

//创建用户
func Login(w http.ResponseWriter, r *http.Request, args httprouter.Params) {
	io.WriteString(w, args.ByName("name"))
}
