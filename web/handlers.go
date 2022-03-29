package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type HomePage struct {
	Name string
}
type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		p := &HomePage{
			Name: "YOURS VIDEO",
		}
		t, err := template.ParseFiles("../templates/home.html")
		if err != nil {
			log.Println("Parsing template home.html error: ", err)
		}
		t.Execute(w, p)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
	}

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.PostFormValue("username")

	var p *UserPage
	if len(cname.Value) != 0 {
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}

	t, e := template.ParseFiles("../templates/userhome.html")
	if e != nil {
		log.Println("parse userhome error ", e)
	}
	t.Execute(w, p)
}

//api透传
func apiHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusNotAcceptable, "Request not recognized")
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, "Request Not recognized")
		return
	}

	request(apibody, w, r)
	defer r.Body.Close()
}

//proxy 转发
func proxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u, _ := url.Parse(proxyURL)
	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.ServeHTTP(w, r)
}
