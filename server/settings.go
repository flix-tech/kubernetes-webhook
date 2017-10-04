package main

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type SettingsCredential struct{
	Key        string `yaml:"key"`
	UserGlobs  []string `yaml:"userglobs"`
	GroupGlobs []string `yaml:"groupglobs"`
}

type Settings struct{
	Credentials *[]SettingsCredential `yaml:"credentials"`
	ListenAddress string `yaml:"listenaddress"`
	SSL        bool `json:"ssl"`
	SSLKeyPath string `json:"sslkeypath"`
	SSLCrtPath string `json:"sslcrtpath"`
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