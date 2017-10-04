package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/prometheus/common/log"
	"github.com/gobwas/glob"
	"io/ioutil"
)

func userMatchesGlobs(username string,globs []string) bool{
	var g glob.Glob
	for _, userGlob := range globs{
		g = glob.MustCompile(userGlob)
		if g.Match(username) {
			return true
		}
	}
	return false
}

func filterValidGroups(groups []string,globs[]string) []string{
	var g glob.Glob
	allowedGroups := []string{}
	for _,group := range groups {
		for _, groupGlob := range globs {
			g = glob.MustCompile(groupGlob)
			if g.Match(group){
				allowedGroups = append(allowedGroups,group)
				break
			}
		}
	}
	return allowedGroups
}
func verifyToken(tokenString string,credentials *[]SettingsCredential) (err1 error, user string, groups []string) {
	var token *jwt.Token = nil
	var err error = nil
	// Iterate over all configured credentials and use the first one that works.
	for _, credential := range(*credentials) {
		publicKey, err := ioutil.ReadFile(credential.Key)
		if err != nil {
			log.Error("RSA key file not found: ", credential.Key, err)
			continue
		}
		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwt.ParseRSAPublicKeyFromPEM(publicKey)
		})
		if err != nil {
			log.Info(err)
			log.Debug("Tried credential: ", string(publicKey))
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
			if userMatchesGlobs(user,credential.UserGlobs){
				log.Info("Success for user", user)
				// ... but allow returning only a subset of groups
				return nil, user, filterValidGroups(groups,credential.GroupGlobs)
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

