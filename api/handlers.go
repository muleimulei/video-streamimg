package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/session"

	"github.com/julienschmidt/httprouter"
)

//创建用户
func CreateUser(w http.ResponseWriter, r *http.Request, args httprouter.Params) {
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	user := &defs.UserCredential{}

	if err := json.Unmarshal(b, user); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if len(user.Username) == 0 || len(user.Pwd) == 0 {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	if err := dbops.AddUserCredential(user.Username, user.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	//加入数据库成功
	id := session.GenerateNewSessionId(user.Username)
	su := &defs.SignedUp{Success: true, SessionId: id}

	data, err := json.Marshal(su)
	if err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	}
	sendNormalResponse(w, string(data), 200)
}

//用户登录
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	b, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	user := &defs.UserCredential{}

	if err := json.Unmarshal(b, user); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	// 验证请求体
	uname := ps.ByName("username")
	// log.Println("login url name : ", uname)
	// log.Println("login body name : ", user.Username)

	if uname != user.Username {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	pwd, err := dbops.GetUserCredential(user.Username)
	log.Println("Login pwd: ", pwd)
	if err != nil || len(pwd) == 0 || pwd != user.Pwd {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	id := session.GenerateNewSessionId(user.Username)
	si := defs.SignedIn{Success: true, SessionId: id}

	if resp, err := json.Marshal(si); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	uname := ps.ByName("username")
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Println("Error in GetUserInfo , ", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	ui := defs.UserInfo{Id: u}

	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	videoBody := &defs.NewVideo{}

	if err := json.Unmarshal(res, videoBody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddNewVideo(videoBody.AutorId, videoBody.Name)
	log.Println("Author Id: ", videoBody.AutorId, " name: ", videoBody.Name)

	if err != nil {
		log.Println("Error in AddNewvideo :", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}

func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Println("listvideos ValidateUser: %n")
		return
	}

	uname := p.ByName("username")
	log.Printf("listvideos url name: %s", uname)

	vs, err := dbops.ListAllVideos(uname, 0, time.Now().UnixMilli())
	if err != nil {
		log.Printf("Error in ListAllVideos: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	vsi := &defs.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	vid := ps.ByName("vid-id")
	if err := dbops.DeleteVideoInfo(vid); err != nil {
		log.Println("Error in DeleteVideo: ", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}

	sendNormalResponse(w, "Delete success", http.StatusNoContent)
}

func PostComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	cbody := &defs.NewComment{}

	if err := json.Unmarshal(reqBody, cbody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := ps.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Println("Error in postcomment : ", err)
		return
	}

	sendNormalResponse(w, "ok", http.StatusCreated)
}

func ShowComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if !ValidateUser(w, r) {
		log.Println("Unathorized user")
		return
	}

	vid := ps.ByName("vid-id")

	cs, err := dbops.ListComments(vid, 0, time.Now().UnixMilli())
	if err != nil {
		log.Println("Error in ListComments ", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}

	sms := &defs.Comments{Comments: cs}
	if data, err := json.Marshal(sms); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(data), http.StatusOK)
	}
}
