package main

const (
	proxyURL = "http://127.0.0.1:9000/"
)

type ApiBody struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	ReqBody string `json:"req_body"`
}
