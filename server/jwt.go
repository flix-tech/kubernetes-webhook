package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/prometheus/common/log"
)

func verifyToken(tokenString string,credentials *[]SettingsCredential) (err1 error, user string, groups []string) {
	var token *jwt.Token = nil
	var err error = nil
	// Iterate over all configured credentials and use the first one that works.
	for _, credential := range(*credentials) {
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
			// Cast the groups to string TODO There has to be a better way
			groups := make([]string,len(claims["groups"].([]interface{})))
			for i, group := range claims["groups"].([]interface{}) {
				groups[i] = group.(string)
			}
			// Do not allow a user mismatch...
			if credential.checkUser(user){
				log.Info("Success for user", user)
				// ... but allow returning only a subset of groups
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

