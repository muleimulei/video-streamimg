package session

import (
	"log"
	"sync"
	"time"
	"video_server/api/dbops"
	"video_server/api/defs"

	uuid "github.com/satori/go.uuid"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func LoadSessionFromDB() {
	r, err := dbops.RetriveAllSession()
	if err != nil {
		log.Println("Rereive session error, ", err)
		return
	}

	r.Range(func(key interface{}, value interface{}) bool {
		ss := value.(*defs.SimpleSession)
		sessionMap.Store(key, ss)
		return true
	})
}

func GenerateNewSessionId(un string) string {
	id, _ := uuid.NewV4()

	ttl := time.Now().Add(5 * time.Minute).UnixMilli()

	ss := &defs.SimpleSession{Username: un, TTL: ttl}

	sessionMap.Store(id.String(), ss)

	dbops.InsertSession(id.String(), ttl, un)
	return id.String()
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		now := time.Now().UnixMilli()

		if session := ss.(*defs.SimpleSession); session.TTL < now { //过期
			dbops.DeleteSessionById(sid)
			sessionMap.Delete(sid)
			return "", true
		}
		return ss.(*defs.SimpleSession).Username, false
	}
	return "", true
}
