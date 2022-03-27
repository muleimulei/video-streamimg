package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
	"video_server/api/defs"
)

func InsertSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)

	stmt, err := db.Prepare("INSERT INTO sessions (sessionid, TTL, loginname) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}
	return nil
}

func RetriveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmt, err := db.Prepare("SELECT TTL, loginname FROM sessions WHERE sessionid = ? ")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var ttl string
	var uname string

	err = stmt.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if ttlint, err := strconv.ParseInt(ttl, 10, 64); err != nil {
		ss.TTL = ttlint
		ss.Username = uname
	} else {
		return nil, err
	}
	return ss, nil
}

func RetriveAllSession() (*sync.Map, error) {
	var m = &sync.Map{}
	stmt, err := db.Prepare(`SELECT * FROM sessions`) // (]

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, ttlstr, loginname string
		if err := rows.Scan(&id, &ttlstr, &loginname); err != nil {
			log.Println("retrive sessions errors: ", err)
			break
		}

		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err != nil {
			ss := &defs.SimpleSession{Username: loginname, TTL: ttl}
			m.Store(id, ss)
		}
	}
	return m, nil
}

func DeleteSessionById(sid string) error {
	stmt, err := db.Prepare("DELETE FROM sessions where sessionid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sid)
	if err != nil {
		return err
	}
	return nil
}
