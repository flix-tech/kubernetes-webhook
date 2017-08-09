package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/gobwas/glob"
)

type SettingsCredential struct{
	Key string
	UserGlob []string
	GroupGlob []string
}

func (c *SettingsCredential) checkUser(user string) bool{
	var g glob.Glob
	for _, userGlob := range c.UserGlob{
		g = glob.MustCompile(userGlob)
		if g.Match(user) {
			return true
		}
	}
	return false
}

func (c *SettingsCredential) checkGroups(groups []string) []string{
	var g glob.Glob
	allowedGroups := []string{}
	for _,group := range groups {
		for _, groupGlob := range c.GroupGlob {
			g = glob.MustCompile(groupGlob)
			if g.Match(group){
				allowedGroups = append(allowedGroups,group)
				break
			}
		}
	}
	return allowedGroups
}

type Settings struct{
	Credentials *[]SettingsCredential
	ListenAddress string
}

func readConfig(path string) (error, *Settings){
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err, nil
	}

	settings := Settings{}
	err = yaml.Unmarshal(data, &settings)
	if err != nil {
		return err, nil
	}

	return nil, &settings
}