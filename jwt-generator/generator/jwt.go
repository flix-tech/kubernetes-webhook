package generator

import (
	"time"
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(iat time.Time, validity time.Duration,user string, groups []string,key *rsa.PrivateKey ) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": iat,
		"exp": iat.Add(validity).Unix(),
		"user": user,
		"groups": groups,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
