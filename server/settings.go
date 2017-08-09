package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type SettingsCredential struct{
	Key        string
	UserGlobs  []string
	GroupGlobs []string
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