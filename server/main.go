package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main() {
	filename := flag.String("config", "config.yml", "Path to config file")
	generateEmptyConfig := flag.Bool("init", false, "Writes an empty config file to the specified path")
	flag.Parse()
	if *generateEmptyConfig {
		out, err := yaml.Marshal(
			Settings{
				Credentials:
				&[]SettingsCredential{{}},
			})
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(*filename,out,0666)
		return
	}
	err, config := readConfig(*filename)
	if err != nil {
		panic(err)
	}
	runServer(config)
}