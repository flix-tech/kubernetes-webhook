package main

import (
	"testing"
	generator "github.com/flix-tech/kubernetes-webhook/jwt-generator"
	"time"
	"io/ioutil"
	"log"
	"github.com/dgrijalva/jwt-go"
	"crypto/rsa"
)

func loadPrivateKey() *rsa.PrivateKey {
	dat, err := ioutil.ReadFile("./RS256.test.key")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(dat)
	if err != nil {
		log.Fatal("Couldn't load private key: ", err)
	}
	return privateKey
}

func loadPublicKey() string {
	dat, err := ioutil.ReadFile("./RS256.test.key.pub")
	if err != nil {
		log.Fatal("Couldn't load public key: ", err)
	}
	return string(dat)
}

func TestHappyCase(t *testing.T) {
	user, groups := "temper", []string{"temper1", "temper2"}
	privKey := loadPrivateKey()
	pubKey := loadPublicKey()
	token, err := generator.GenerateToken(time.Now().Add(-1* time.Second),60 * time.Second,user, groups,privKey)
	if err != nil {
		log.Fatal("Generating token failed", err)
	}

	err1, user1, _ := verifyToken(token,&[]SettingsCredential{{Key: pubKey,UserGlob:"temper",GroupGlob:"temper*"}})
	if err1 != nil {
		t.Error(err1)
	}
	if user1 != user {
		t.Error("Users did not match")
	}
}
