package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/gobwas/glob"
)

type SettingsCredential struct{
	Key string
	UserGlob string
	GroupGlob string
}

func (c *SettingsCredential) checkUser(user string) bool{
	var g glob.Glob
	g = glob.MustCompile(c.UserGlob)
	return g.Match(user)
}

func (c *SettingsCredential) checkGroups(groups []string) []string{
	var g glob.Glob
	g = glob.MustCompile(c.GroupGlob)
	allowedGroups := []string{}
	for _,group := range groups {
		if g.Match(group){
			allowedGroups = append(allowedGroups,group)
		}
	}
	return allowedGroups
}

type Settings struct{
	Credentials *[]SettingsCredential
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