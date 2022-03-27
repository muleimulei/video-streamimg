package defs

import "net/http"

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrorResponse struct {
	HttpSc int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{HttpSc: http.StatusBadRequest, Error: Err{Error: "Request body is not correct", ErrorCode: "001"}}
	ErrorNotAuthUser            = ErrorResponse{HttpSc: http.StatusUnauthorized, Error: Err{Error: "User anthentication failed.", ErrorCode: "002"}}
	ErrorDBError                = ErrorResponse{HttpSc: http.StatusInternalServerError, Error: Err{Error: "DB ops failed", ErrorCode: "003"}}
	ErrorInternalFaults         = ErrorResponse{HttpSc: http.StatusInternalServerError, Error: Err{Error: "Internal service error", ErrorCode: "004"}}
)
