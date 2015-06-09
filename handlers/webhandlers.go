package handlers

import (
	"log"

	"../companies"
	"../memberships"
	"../people"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetPeople(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	people, err := people.GetAll(db)
	if err != nil {
		log.Printf("GetPeople: %s", err)
		return
	}
	c.JSON(200, people)
}

func GetCompanies(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	cos, err := companies.GetAll(db)
	if err != nil {
		log.Printf("GetCompanies: %s", err)
		return
	}
	c.JSON(200, cos)
}

func GetMemberships(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	mem, err := memberships.GetAll(db)
	if err != nil {
		log.Printf("GetMemberships: %s", err)
		return
	}
	c.JSON(200, mem)
}

func GetCompanyById(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	cid := c.Param("cid")
	//Add POC here
	cos, err := companies.GetById(db, cid)
	if err != nil {
		log.Printf("GetCompanyById: %s", err)
		return
	}
	c.JSON(200, cos)
}

func GetPersonById(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	pid := c.Param("pid")
	//Add POC here
	p, err := people.GetById(db, pid)
	if err != nil {
		log.Printf("GetPersonById: %s", err)
		return
	}
	c.JSON(200, p)
}
