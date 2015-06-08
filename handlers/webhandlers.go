package handlers

import "github.com/gin-gonic/gin"

func GetPeople(c *gin.Context) {
	c.JSON(200, "People")
}
