package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + "/" + vid

	video, err := os.Open(vl)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
}

//上传
func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}
	file, _, err := r.FormFile("file") //

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal error")
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Read file error : ", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal error")
		return
	}

	fn := p.ByName("vid-id")
	log.Println("video id :", fn)
	err = ioutil.WriteFile(VIDEO_DIR+"/"+fn, data, 0644)
	if err != nil {
		log.Println("Write file error : ", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal error")
		return
	}

	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload finished")
}
