package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)


func logFatalError(sector string, err error) {
	if err != nil {
		log.Println("=>", sector)
		log.Fatal("=> ", err)
	}
}
func main() {
	// sql.DB 객체 생성
	db, err := sql.Open("mysql", "root:wowplus@tcp(127.0.0.1:3306)/official_api")
	logFatalError("Connection", err)

	defer db.Close()

	// 트랜젝션 시작
	tx, err := db.Begin()
	logFatalError("TransBegin", err)

	defer tx.Rollback()	// 트랜젝션 중간에 에러발생 시 Rollback 처리

	// INSERT 문 실행
	_, err = db.Exec("INSERT INTO dbsessions VALUES (?, ?, ?, UNIX_TIMESTAMP(), ?)", "sample_insert_by_golang4", "127.0.0.1", "IntelliJ IDEA 2016.2.5-community_version", "{no data}")
	logFatalError("Execute 1st Query Statement", err)

	_, err = db.Exec("INSERT INTO dbsessions VALUES (?, ?, ?, UNIX_TIMESTAMP(), ?)", "sample_insert_by_golang5", "127.0.0.1", "IntelliJ IDEA 2016.2.5-community_version", "{no data}")
	logFatalError("Execute 2nd Query Statement", err)

	// 트랜젝션 Commit 처리
	err = tx.Commit()
	logFatalError("TransCommit", err)
}
