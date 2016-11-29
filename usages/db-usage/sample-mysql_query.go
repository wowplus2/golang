package main
// Mysql 사용 - 쿼리 사용하기
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)

// placeholder:
// 	MySql -> ?, ?, ...
//	Oracle -> :val1, :val2, ...
//	PostgreSQL -> $1, $2, ...
//	MS-SQL -> ?, ?n, :n, $n, ... 등을 지원

func mysql_connection() {
	// sql.DB 객체 생성
	// db, err := sql.Open("mysql", "user:password@/dbname")
	// db, err := sql.Open("mssql", "server=(local);user id=sa;password=pwd;database=dbname")
	db, err := sql.Open("mysql", "root:wowplus@tcp(127.0.0.1:3306)/gnu5board")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
}

// 하나의 Row를 갖는 SQL Query
func mysql_singlerow() {
	// sql.DB 객체 생성
	//  db, err := sql.Open("mysql", "user:password@/dbname")
	db, err := sql.Open("mysql", "root:wowplus@tcp(127.0.0.1:3306)/gnu5board")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var name string
	err = db.QueryRow("SELECT bo_table FROM g5_board WHERE gr_id = ?", "genernal").Scan(&name)
	/*if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\t",name)*/
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID...")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Println("\t",name)
	}
}

// 복수의 Row를 갖는 SQL Query
func mysql_multirows() {
	// sql.DB 객체 생성
	//  db, err := sql.Open("mysql", "user:password@/dbname")
	db, err := sql.Open("mysql", "root:wowplus@tcp(127.0.0.1:3306)/official_api")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var act int
	var seid string
	rows, err := db.Query("SELECT session_id, last_activity FROM dbsessions ORDER BY last_activity DESC")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()	// 반드시 닫기!!(지연하여 닫는다.)

	for rows.Next() {
		err := rows.Scan(&seid, &act)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\tsession_id :",seid, ", last_activity :", act)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	fmt.Println("mysql_singlerow: ")
	mysql_singlerow()
	fmt.Println()
	fmt.Println("mysql_multirows: ")
	mysql_multirows()
}
