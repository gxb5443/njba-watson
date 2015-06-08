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

func GetCompanies(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	query := `
		SELECT id, name FROM companies;	
	`
	var cos []struct {
		Id   string
		Name string `db:"name"`
	}
	err := db.Select(&cos, query)
	if err != nil {
		log.Printf("GetCompanies: %s", err)
		return
	}
	c.JSON(200, cos)
}

func GetMemberships(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	query := `
		SELECT companies.name as co, companies.id as co_id, membership_holdings.member_id as mem from membership_holdings inner join companies on membership_holdings.company_id=companies.id;
	`
	var cos []struct {
		MemId string `db:"mem"`
		CoId  string `db:"co_id"`
		Name  string `db:"co"`
	}
	err := db.Select(&cos, query)
	if err != nil {
		log.Printf("GetMemberships: %s", err)
		return
	}
	c.JSON(200, cos)
}
