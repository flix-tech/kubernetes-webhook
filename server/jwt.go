package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/prometheus/common/log"
)

func verifyToken(tokenString string,credentials *[]SettingsCredential) (err1 error, user string, groups []string) {
	var token *jwt.Token = nil
	var err error = nil
	for _, credential := range(*credentials) {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return jwt.ParseRSAPublicKeyFromPEM([]byte(credential.Key))
		})
		if err != nil {
			log.Info(err)
			log.Debug("Tried credential: ", credential.Key)
			continue
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Expiration checking
			user := claims["user"].(string)
			groups := make([]string,len(claims["groups"].([]interface{})))
			for i, group := range claims["groups"].([]interface{}) {
				groups[i] = group.(string)
			}
			//groups := claims["groups"].([]string)
			if credential.checkUser(user){
				// return the user and the allowed groups
				log.Info("Success for user", user)
				return nil, user, credential.checkGroups(groups)
			}else {
				err = errors.New("User not allowed by glob")
			}
		}
	}
	if err != nil {
		return errors.New("Could not find a matching token"), "", nil
	}
	panic("Impossible codepath reached")
}

