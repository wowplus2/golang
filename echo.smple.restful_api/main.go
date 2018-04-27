package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello for the web side!")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	dataType := c.Param("data")

	switch dataType {
	case "string":
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is %s\nand his type is %s\n", catName, catType))
	case "json":
		return c.JSON(http.StatusOK, map[string]map[string]string{
			"result": {
				"resp":    "true",
				"code":    "200",
				"message": "ok",
			},
			"data": {
				"name": catName,
				"type": catType,
			},
		})
	default:
		return c.JSON(http.StatusOK, map[string]map[string]string{
			"result": {
				"resp":    "true",
				"code":    "200",
				"message": "ok",
			},
			"data": {
				"name": catName,
				"type": catType,
			},
		})
	}
}

// ioutil.ReadAll() is fastest in 3way
func addCats(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addCat: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling in addCat: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your cat: %#v", cat)
	return c.String(http.StatusOK, "We got your cat!")
}

// json.NewDecoder has normal speed in 3way(got it!)
func addDogs(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed processing addDog request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your dog %#v", dog)
	return c.String(http.StatusOK, "we got your dog!")
}

// Bind() is slowest in 3way
func addHansters(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Failed processing addHamster request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your hamster %#v", hamster)
	return c.String(http.StatusOK, "we got your hamster!")
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	e.GET("/", hello)
	e.GET("/cats/:data", getCats)

	e.POST("/cats", addCats)
	e.POST("/dogs", addDogs)
	e.POST("/hamsters", addHansters)

	e.Start(":3000")
}
