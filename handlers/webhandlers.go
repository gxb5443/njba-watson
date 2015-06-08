package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetPeople(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	query := `
		SELECT id, first_name, last_name FROM people;	
	`
	var people []struct {
		Id    string
		Fname string `db:"first_name"`
		Lname string `db:"last_name"`
	}
	err := db.Select(&people, query)
	if err != nil {
		log.Printf("GetPeople: %s", err)
		return
	}
	c.JSON(200, people)
}
