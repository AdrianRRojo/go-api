package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateJWT creates a JWT with custom claims and signs it using a secret key.
func createJWT(companyID, email string, exp int64) (string, error) {
	claims := jwt.MapClaims{
		"companyID": companyID,
		"email":     email,
		"exp":       exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("API_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// generateJWT creates the tokenStruct and inserts it
func generateJWT(collection *mongo.Collection) {
	companyID := os.Getenv("COMPANY")
	email := os.Getenv("COMPANY_EMAIL")
	exp := time.Now().Add(3 * 365 * 24 * time.Hour).Unix()

	token, err := createJWT(companyID, email, exp)
	if err != nil {
		log.Fatal(err)
	}

	newCompany := tokenStruct{
		CompanyID: companyID,
		Email:     email,
		Exp:       exp,
		Token:     token,
	}

	id := insertNewCompany(collection, newCompany)

	fmt.Println("Inserted Document ID:", id)
	fmt.Println("Generated JWT:", token)
}
