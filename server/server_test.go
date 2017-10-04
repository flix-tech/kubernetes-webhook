package main

import (
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	config := Settings{
		SSL:true,
		SSLCrtPath: "./test_resources/cert.pem",
		SSLKeyPath: "./test_resources/key.pem",
		Credentials: &[]SettingsCredential{
			{
				Key: "./test_resources/RS256.test.pub",
				GroupGlobs: []string{"*"},
				UserGlobs: []string{"*"},
			},
		},
		ListenAddress: "127.0.0.1:55368",
	}
	print("Starting server...")
	srv := runServer(&config)
	println("Server runs")
	time.Sleep(1* time.Second)
	println("Server timeout exceeded")
	if err := srv.Shutdown(nil); err != nil {
		t.Error("Shutting down server failed", err)
	}
}
