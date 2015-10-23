package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	rando "math/rand"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//GenerateRandomBytes returns securely generated, random bytes.
//It returns an error if the system's securet random number gen
//fails, in which everything should stop working
func GenerateRandomBytes(n int64) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	//Err==nill iff we read len(b) bytes
	if err != nil {
		return nil, err
	}
	return b, nil
}

//GenerateRandomString uses generate random bytes to generate a
//URL-safe, base64 encoded string.
func GenerateRandomString(n int64) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

//GenerateJWT generates Javascript Web Token for
func GenerateJWT(u interface{}, claim_name string) (string, error) {
	RSA_KEY := os.Getenv("RSA_KEY")
	if RSA_KEY == "" {
		panic("Utils.GenerateJWT: No key specified")
	}
	JWT_LIFE := os.Getenv("JWT_LIFE")
	if JWT_LIFE == "" {
		JWT_LIFE = "3600"
	}
	jwt_life, err := strconv.Atoi(JWT_LIFE)
	if err != nil {
		panic("Utils.GenerateJWT: JWT_LIFE must be an integer numeber")
	}

	alg := jwt.GetSigningMethod("HS256")
	token := jwt.New(alg)
	token.Header["typ"] = "JWT"
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Second * time.Duration(jwt_life)).Unix()
	token.Claims[claim_name] = &u
	if out, terr := token.SignedString([]byte(RSA_KEY)); terr == nil {
		return out, nil
	} else {
		return "", terr
	}
}

//Quickly authenticates provided JWT
func ParseJWT(jwt_token string) (map[string]interface{}, bool, error) {
	token, err := jwt.Parse(jwt_token, func(t *jwt.Token) (interface{}, error) {
		RSA_KEY := os.Getenv("RSA_KEY")
		if RSA_KEY == "" {
			return nil, errors.New("Utils.ParseJWT: No key specified")
		}
		return []byte(RSA_KEY), nil
	})
	return token.Claims, token.Valid, err
}

func GeneratePassword(n int64) string {
	rando.Seed(time.Now().Unix())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	available_chars := len(chars)
	password := make([]rune, n)
	for i := range password {
		password[i] = chars[rando.Intn(available_chars)]
	}
	return string(password)
}
