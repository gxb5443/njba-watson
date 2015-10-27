package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/gxb5443/Cuddy/tokens"
	"github.com/gxb5443/Cuddy/users"

	"github.com/gxb5443/Cuddy/credentials"
	"github.com/gxb5443/Cuddy/utils"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

//Login Handler
func Login(c *gin.Context) {
	var loginreq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	db := c.MustGet("db").(*sqlx.DB)
	c.Bind(&loginreq)
	if loginreq.Username == "" || loginreq.Password == "" {
		c.JSON(401, gin.H{"status": "Username and Password Required"})
		log.Println("Blank Credentials")
		return
	}
	u, err := credentials.Login(db, loginreq.Username, loginreq.Password)
	if err != nil {
		//c.JSON(200, gin.H{"status": "Invalid Username/Password"})
		c.AbortWithError(404, errors.New("Username/Password combination not found"))
		log.Println(err)
		return
	}
	ju, jerr := json.Marshal(u)
	if jerr != nil {
		//c.JSON(200, gin.H{"status": "Error encoding Credentials"})
		c.AbortWithError(500, errors.New("Error encoding Credentials"))
		log.Println(jerr)
		return
	}

	token, terr := utils.GenerateJWT(string(ju), "usr")
	if terr != nil {
		//c.JSON(200, gin.H{"status": "Error encoding Credentials"})
		c.AbortWithError(500, errors.New("Error generating JWT"))
		log.Println(terr)
		return
	}

	//Fetch refresh token
	refresh, rerr := tokens.GetByUserId(db, u.Id)
	if rerr != nil && rerr != tokens.ErrNoTokensFound {
		c.AbortWithError(500, errors.New("Error getting Refresh Token"))
		log.Println(rerr)
		return
	}
	//if len(refresh) == 0 {
	if refresh == nil {
		c.JSON(200, gin.H{"token": token, "refresh": ""})
		return
	}
	//c.JSON(200, gin.H{"token": token, "refresh": refresh[0].Id})
	c.JSON(200, gin.H{"token": token, "refresh": refresh.Token})
}

func Logout(c *gin.Context) {
	//session := c.MustGet("session_manager").(*sessions.Manager)
	//session.Destroy(c.Writer, c.Request)
}

func RefreshToken(c *gin.Context) {
	db := c.MustGet("db").(*sqlx.DB)
	var refresh struct {
		Token string `json:"token"`
	}
	c.Bind(&refresh)
	uid, err := tokens.Verify(db, refresh.Token)
	if err == tokens.ErrTokenNotValid {
		c.AbortWithError(404, errors.New("Invalid Token"))
		log.Println(err)
		return
	}
	if err != nil {
		log.Println(err)
		c.AbortWithError(500, errors.New("Error Verifying Token"))
		return
	}
	var user users.User
	log.Println(reflect.TypeOf(user))
	u, uerr := users.Get(db, uid)
	if uerr != nil {
		log.Println(uerr)
		c.AbortWithError(500, errors.New("Could not get user"))
		return
	}
	ju, jerr := json.Marshal(u)
	if jerr != nil {
		log.Println(jerr)
		c.AbortWithError(500, errors.New("Error encoding Credentials"))
		return
	}
	token, terr := utils.GenerateJWT(string(ju), "usr")
	if terr != nil {
		log.Println(terr)
		c.AbortWithError(500, errors.New("Error generating JWT"))
		return
	}
	c.JSON(200, gin.H{"token": token})
	return
}
