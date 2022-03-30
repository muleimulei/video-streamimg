package dbops

import (
	"database/sql"
	"log"
	"time"
	"video_server/api/defs"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

var db *sql.DB

//连接数据库
func init() {
	var err error
	str := `mulei:Mulei666/@tcp(8.142.31.201:3306)/videoserver?charset=utf8`
	db, err = sql.Open("mysql", str)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func AddUserCredential(loginname string, pwd string) error {
	stmt, err := db.Prepare("INSERT INTO users (login_name, pwd) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(loginname, pwd)
	if err != nil {
		return err
	}
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmt, err := db.Prepare("SELECT pwd FROM users WHERE login_name = ?")
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer stmt.Close()

	var pwd sql.NullString
	err = stmt.QueryRow(loginName).Scan(&pwd)
	if err != nil {
		return "", err
	}
	return pwd.String, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmt, err := db.Prepare("DELETE FROM users where login_name = ? and pwd = ?")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	return nil
}

func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	ctime := time.Now().Format("Jan 02 2006, 15:04:05")

	stmt, err := db.Prepare("INSERT INTO videoinfo (id, authorid, name, display_ctime) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id.String(), aid, name, ctime)
	if err != nil {
		return nil, err
	}
	res := &defs.VideoInfo{
		Id:           id.String(),
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: ctime,
	}
	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmt, err := db.Prepare("SELECT authorid, name, display_ctime where id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var aid int
	var name string
	var dct string

	err = stmt.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &defs.VideoInfo{
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: dct,
	}
	return res, nil
}

func DeleteVideoInfo(vid string) error {
	stmt, err := db.Prepare("DELETE FROM videoinfo where id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(vid)
	if err != nil {
		return err
	}
	return nil
}

func AddNewComments(vid string, aid int, content string) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO comments (id, videoid, authorid, content) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id.String(), vid, aid, content)
	if err != nil {
		return err
	}
	return nil
}

// // where comments.videoid = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)
func ListComments(vid string, from, to int64) ([]*defs.Comment, error) {
	stmt, err := db.Prepare(`SELECT comments.id, users.login_name, comments.content FROM comments INNER JOIN users
	ON comments.authorid = users.id  where comments.videoid = ? order by time desc`)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var res []*defs.Comment

	rows, err := stmt.Query(vid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}
		c := &defs.Comment{
			Id:         id,
			VideoId:    vid,
			AuthorName: name,
			Content:    content,
		}
		res = append(res, c)
	}
	log.Println(vid, " ", res)
	return res, nil
}

func GetUser(uname string) (int, error) {
	stmt, err := db.Prepare("SELECT id FROM users where login_name = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	var uid int
	err = stmt.QueryRow(uname).Scan(&uid)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

func ListAllVideos(uname string, from, to int64) ([]*defs.VideoInfo, error) {
	stmt, err := db.Prepare(`SELECT videoinfo.id, videoinfo.authorid, videoinfo.name, videoinfo.display_ctime FROM videoinfo
		INNER JOIN users on videoinfo.authorid = users.id WHERE users.login_name = ?`)

	if err != nil {
		log.Println("ListAllVideos error , ", err)
		return nil, err
	}

	defer stmt.Close()
	var res []*defs.VideoInfo

	rows, err := stmt.Query(uname)
	if err != nil {
		log.Println("ListAllVideos error , ", err)
		return nil, err
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int

		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}

		vi := &defs.VideoInfo{
			Id:           id,
			AuthorId:     aid,
			Name:         name,
			DisplayCtime: ctime,
		}
		res = append(res, vi)
	}
	return res, nil
}
