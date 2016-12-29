package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:wowplus@tcp(127.0.0.1:3306)/gotest")
	if err != nil {
		fmt.Println(err.Error())
	}

	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	stmt, err := db.Prepare("CREATE TABLE person (idx int NOT NULL AUTO_INCREMENT, first_name varchar(40), last_name varchar(40), PRIMARY KEY (idx));")
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Person Table successfully migration...")
	}

}
