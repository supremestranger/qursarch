package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	SIGNING_METHOD = jwt.SigningMethodHS256
)

// Function to create JWT tokens with claims
func CreateToken(username string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(SIGNING_METHOD, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)

	tokenString, err := claims.SignedString("SECRET_KEY")
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenStr string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) { return "SECRET_KEY", nil })
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	// if claims["iat"]

	if token.Method.Alg() != SIGNING_METHOD.Name {
		return false
	}

	return true
}
