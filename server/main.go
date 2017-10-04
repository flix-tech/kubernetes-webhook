package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"fmt"
	"os"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	filename := flag.String("config", "config.yml", "Path to config file")
	generateEmptyConfig := flag.Bool("init", false, "Writes an empty config file to the specified path")
	flag.Parse()
	if *generateEmptyConfig {
		filename = createEmptyConfig(filename)
		return
	}
	err, config := readConfig(*filename)
	if err != nil {
		panic(err)
	}
	if len(*config.Credentials) == 0{
		log.Fatal("No credentials to verify were provided")
	}
	srv := runServer(config)
	<-done
	if err := srv.Shutdown(nil); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func createEmptyConfig(filename *string) *string {
	out, err := yaml.Marshal(
		Settings{
			Credentials:
			&[]SettingsCredential{{}},
		})
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(*filename, out, 0666)
	return filename
}