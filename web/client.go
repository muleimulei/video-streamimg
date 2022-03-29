package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	switch b.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", b.Url, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println("api proxy request ", err)
			return
		}
		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", b.Url, strings.NewReader(b.ReqBody))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println("api proxy request ", err)
			return
		}
		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("Delete", b.Url, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println("api proxy request ", err)
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
	}
}

func normalResponse(w http.ResponseWriter, r *http.Response) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(b))

}
