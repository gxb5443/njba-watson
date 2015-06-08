package main

import (
	"database/sql"
	"devportal/applications"
	"devportal/credentials"
	"devportal/users"
	"devportal/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"../Cuddy/handlers"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var u = users.User{
	FirstName: "Minnie",
	LastName:  "Strator",
	Email:     "admin@namely.com",
	Admin:     true,
}

var sa = applications.Application{
	Name: "Default 1",
}

var login = credentials.Credential{
	PasswordHash: "password",
}

func main() {
	var DBOPEN string
	LISTEN_WEB := os.Getenv("LISTEN_WEB")
	if LISTEN_WEB == "" {
		LISTEN_WEB = ":8080"
	}
	LISTEN_API := os.Getenv("LISTEN_API")
	if LISTEN_API == "" {
		LISTEN_API = ":8000"
	}
	SSLMODE := os.Getenv("SSLMODE")
	if SSLMODE == "" {
		SSLMODE = "disable"
	}
	DBOPEN += fmt.Sprintf("sslmode=" + SSLMODE)
	USER := os.Getenv("USER")
	if USER == "" {
		log.Fatalln("USER ENV Variable not set")
	}
	DBOPEN += fmt.Sprint(" user=" + USER)
	DBPASS := os.Getenv("DBPASS")
	if DBPASS != "" {
		DBOPEN += fmt.Sprint(" password=" + DBPASS)
	}
	DBNAME := os.Getenv("DBNAME")
	if DBNAME == "" {
		log.Fatalln("DBNAME ENV Variable not set")
	}
	DBOPEN += fmt.Sprint(" dbname=" + DBNAME)
	DBHOST := os.Getenv("DBHOST")
	if DBHOST != "" {
		DBOPEN += fmt.Sprint(" host=" + DBHOST)
	}

	//initUser(DBOPEN)
	RunWeb(DBOPEN, LISTEN_WEB)
}

func RunWeb(dbopen, listen string) {
	defer log.Println("Shutting down server...")
	r := gin.Default()
	//r.Use(gzip.Gzip(gzip.BestSpeed))
	r.Use(DB(dbopen))
	r.GET("/people", handlers.GetPeople)
	/*
		r.POST("/login", handlers.Login)
		r.GET("/logout", handlers.Logout)
	*/
	/*
		v1 := r.Group("/v1")
		v1.Use(JWT())
		{
			v1.GET("/users", handlers.GetAllUsers)
			v1.GET("/users/:id", handlers.GetUser)
			v1.GET("/users/:id/apps", handlers.GetUserApps)
			v1.GET("/apps", handlers.GetAllApps)
			v1.GET("/apps/:id", handlers.GetApp)
			v1.POST("/users", handlers.AddUser)
			v1.PUT("/users", handlers.UpdateUser)
			v1.POST("/apps", handlers.AddApp)
			v1.POST("/secret", handlers.RegenerateAppSecret)
			v1.DELETE("/users/:id", handlers.DeleteUser)
			v1.DELETE("/apps/:id", handlers.DeleteApp)
		}
	*/
	r.NoRoute(static.Serve("/", static.LocalFile("./public/", true)))
	log.Println("Running Webserver...")
	r.Run(listen)
}

func DB(dbopen string) gin.HandlerFunc {
	db, err := sql.Open("postgres", dbopen)
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Request.Header["Authorization"]
		if !ok {
			c.Fail(401, errors.New("Unauthorized: No Authorization Header"))
			return
		}
		auth := c.Request.Header["Authorization"][0]
		authorization_header := strings.Split(auth, " ")
		if authorization_header[0] != "Bearer" {
			c.Fail(401, errors.New("Unauthorized: No Bearer Token"))
			return
		}
		claims, valid, err := utils.ParseJWT(authorization_header[1])
		if err != nil {
			c.Fail(401, err)
			return
		}
		if !valid {
			//			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.Fail(401, errors.New("Unauthorized: Token not valid"))
			return
		}
		c.Set("token", claims)
		c.Set("user", claims["usr"])
		c.Next()
	}
}