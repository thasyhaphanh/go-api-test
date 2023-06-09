package main

import (
	"database/sql"
	"time"

	"net/http"

	"github.com/gin-contrib/cors"

	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// album represents data about a record album.
type Attraction struct {
	Id         string `db."id" json:"id"`
	Name       string `db."name" json:"name"`
	Detail     string `db."detail" json:"detail"`
	Coverimage string `db."coverimage" json:"coverimage"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/gomysql")

	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	router := gin.Default()
	router.GET("/attractions", getAttractions)

	router.Use(cors.Default())
	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAttractions(c *gin.Context) {
	var attractions []Attraction
	rows, err := db.Query("select id, name, detail, coverimage from attractions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var a Attraction
		err := rows.Scan(&a.Id, &a.Name, &a.Detail, &a.Coverimage)
		if err != nil {
			log.Fatal(err)
		}
		attractions = append(attractions, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, attractions)
}
