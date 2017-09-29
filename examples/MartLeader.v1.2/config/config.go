package config

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB
var err error

func Dbconn() (db *sql.DB, err error) {
	//host := "192.168.0.100"
	host := "localhost"
	port := "1433"
	dbname := "iampos_825"
	uid := "martreader"
	pass := "trian@akxmflej"

	// sql.DB 객체 생성
	db, err = sql.Open("mssql", "server="+host+";port="+port+";user id="+uid+"; password="+pass+"; database="+dbname+";encrypt=disable")
	return
}
