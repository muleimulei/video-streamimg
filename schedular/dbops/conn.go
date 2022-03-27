package dbops

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
