package main

import (
	"encoding/json"
	"io"
	"net/http"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/session"

	"github.com/julienschmidt/httprouter"
)

//创建用户
func CreateUser(w http.ResponseWriter, r *http.Request, args httprouter.Params) {
	user := r.PostFormValue("user_name")
	pwd := r.PostFormValue("pwd")

	if len(user) == 0 || len(pwd) == 0 {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(user, pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	//加入数据库成功
	id := session.GenerateNewSessionId(user)
	su := &defs.SignedUp{Success: true, SessionId: id}

	data, err := json.Marshal(su)
	if err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	}
	sendNormalResponse(w, string(data), 200)
}

//创建用户
func Login(w http.ResponseWriter, r *http.Request, args httprouter.Params) {
	io.WriteString(w, args.ByName("name"))
}
