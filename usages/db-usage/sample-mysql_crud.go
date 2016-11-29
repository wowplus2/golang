package main
// Mysql 사용 - DML 사용하기
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)

func main() {
	// sql.DB 객체 생성
	db, err := sql.Open("mysql", "root:wowplus@tcp(127.0.0.1:3306)/official_api")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Execute insert query statment
	res, err := db.Exec("INSERT INTO dbsessions VALUES (?, ?, ?, UNIX_TIMESTAMP(), ?)", "sample_insert_by_golang4", "127.0.0.1", "IntelliJ IDEA 2016.2.5", "{no data}")
	if err != nil {
		log.Fatal(err)
	}

	// sql.Result.LastInsertId() 체크
	/*idx, err := res.LastInsertId()
	if idx > 0 {
		fmt.Println("1 row inserted...ID: ", idx)
	}*/
	// sql.Result.RowsAffected() 체크
	n, err := res.RowsAffected()
	if n == 1 {
		fmt.Println("1 row inserted...")
	}
}
