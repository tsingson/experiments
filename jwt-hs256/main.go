package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sanity-io/litter"
)

func main() {
	mySigningKey := []byte("theverylongverylongverylongpassword")
	// Create the Claims
	type MyCustomClaims struct {
		Role  string `json:"role"`
		Email string `json:"email"`
		jwt.StandardClaims
	}

	// Create the Claims
	customClaims := MyCustomClaims{
		"todo_user",
		"tsingson@mycompany.com",
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),
			ExpiresAt: int64(time.Now().Unix() + 1000),
			Issuer:    "test",
		},
	}

	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token1.SignedString(mySigningKey)
	fmt.Println("custom claims: ")
	litter.Dump(tokenString)

	tt1, _ := jwt.Parse(tokenString, func(*jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	// fmt.Println("verify ")
	litter.Dump(tt1.Claims)
	// litter.Dump(mySigningKey)
}
