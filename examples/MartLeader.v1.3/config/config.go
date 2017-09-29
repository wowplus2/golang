package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func Dbconn() (db *sql.DB, err error) {
	//host := "192.168.10.201"
	host := "localhost"
	port := "3306"
	dbname := "myposys_db_833"
	uid := "mssaf"
	pass := "trian@akxmflej"

	// sql.DB 객체 생성
	db, err = sql.Open("mysql", uid+":"+pass+"@tcp("+host+":"+port+")/"+dbname)
	return
}
