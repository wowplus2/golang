package main
// Mysql 사용 - DML 사용하기
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)


func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// sql.DB 객체 생성
	db, err := sql.Open("mysql", "root:wowplus@tcp(127.0.0.1:3306)/official_api")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Prepared Statement 생성
	stmt, err := db.Prepare("UPDATE dbsessions SET user_agent = ?, last_activity = UNIX_TIMESTAMP() WHERE session_id = ?")
	checkError(err)

	defer stmt.Close()

	// Prepared Statement 실행
	_, err = stmt.Exec("IntelliJ IDEA 2016.2.5-community_version", "sample_insert_by_golang1")
	checkError(err)
	_, err = stmt.Exec("IntelliJ IDEA 2016.2.5-community_version", "sample_insert_by_golang2")
	checkError(err)
	_, err = stmt.Exec("IntelliJ IDEA 2016.2.5-community_version", "sample_insert_by_golang3")
	checkError(err)
}
