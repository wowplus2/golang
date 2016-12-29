package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"bytes"
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

	type Person struct {
		Idx	int
		First_Name	string
		Last_Name	string
	}
	router := gin.Default()

	// Add API handlers here

	// GET a person detail
	router.GET("/person/:id", func(c *gin.Context) {
		var (
			person	Person
			result	gin.H
		)

		idx := c.Param("id")
		row := db.QueryRow("SELECT idx, first_name, last_name FROM person WHERE idx = ?;", idx)
		err = row.Scan(&person.Idx, &person.First_Name, &person.Last_Name)
		if err != nil {
			// If no results send null
			result = gin.H{
				"result" : nil,
				"count" : 0,
			}
		} else {
			result = gin.H{
				"result" : person,
				"count" : 1,
			}
		}

		c.JSON(http.StatusOK, result)
	})

	// GET all persons
	router.GET("/persons", func(c *gin.Context) {
		var (
			person	Person
			persons	[]Person
		)

		rows, err := db.Query("SELECT idx, first_name, last_name FROM person ORDER BY idx DESC;")
		if err != nil {
			fmt.Println(err.Error())
		}

		for rows.Next() {
			err = rows.Scan(&person.Idx, &person.First_Name, &person.Last_Name)
			persons = append(persons, person)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		defer db.Close()

		c.JSON(http.StatusOK, gin.H{
			"result" : persons,
			"count" : len(persons),
		})
	})

	// POST new person details
	router.POST("/person", func(c *gin.Context) {
		var buff bytes.Buffer
		first_name := c.PostForm("fst_nm")
		last_name := c.PostForm("lst_nm")

		stmt, err := db.Prepare("INSERT INTO person (first_name, last_name) VALUES(?, ?);")
		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(first_name, last_name)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Fastest way to append strings
		buff.WriteString(first_name)
		buff.WriteString(" ")
		buff.WriteString(last_name)
		defer stmt.Close()

		name := buff.String()
		c.JSON(http.StatusOK, gin.H{
			"message" : fmt.Sprintf(" %s successfully creates.", name),
		})

	})

	// Delete resources
	router.DELETE("/person", func(c *gin.Context) {
		idx := c.Query("id")
		stmt, err := db.Prepare("DELETE FROM person WHERE idx = ?;")
		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(idx)
		if err != nil {
			fmt.Println(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message" : fmt.Sprintf("Successfully deleted user: %s", idx),
		})
	})

	// PUT - update a person details
	router.PUT("/person", func(c *gin.Context) {
		var buff bytes.Buffer
		idx := c.Query("id")
		first_name := c.PostForm("fst_nm")
		last_name := c.PostForm("lst_nm")

		stmt, err := db.Prepare("UPDATE person SET first_name = ?, last_name = ? WHERE idx = ?;")
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = stmt.Exec(first_name, last_name, idx)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Fastest way to append strings
		buff.WriteString(first_name)
		buff.WriteString(" ")
		buff.WriteString(last_name)
		defer stmt.Close()

		name := buff.String()
		c.JSON(http.StatusOK, gin.H{
			"message" : fmt.Sprintf("Successfully updated to %s", name),
		})
	})
	router.Run(":9100")
}
