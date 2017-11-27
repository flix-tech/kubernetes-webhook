package main

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/flix-tech/kubernetes-webhook/jwt-generator/generator"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

type JWTTestCase struct {
	iat       time.Time
	validity  time.Duration
	user      string
	userGlob  string
	groups    []string
	groupGlob string
	t         *testing.T
	pubKey    string
	privKey   *rsa.PrivateKey
}

func newJWTTestCase(iat time.Time, validity time.Duration, user string, groups []string, userGlob string, groupGlob string, t *testing.T) JWTTestCase {
	privKey := loadPrivateKey()
	pubKey := "./test_resources/RS256.test.key.pub"
	return JWTTestCase{
		iat:       iat,
		validity:  validity,
		user:      user,
		userGlob:  userGlob,
		groups:    groups,
		groupGlob: groupGlob,
		t:         t,
		privKey:   privKey,
		pubKey:    pubKey,
	}
}

func loadPrivateKey() *rsa.PrivateKey {
	dat, err := ioutil.ReadFile("./test_resources/RS256.test.key")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(dat)
	if err != nil {
		log.Fatal("Couldn't load private key: ", err)
	}
	return privateKey
}

func loadPublicKey() string {
	dat, err := ioutil.ReadFile("./test_resources/RS256.test.key.pub")
	if err != nil {
		log.Fatal("Couldn't load public key: ", err)
	}
	return string(dat)
}

func (tc *JWTTestCase) createAndCheckToken() (error, string, []string) {
	token, err := generator.GenerateToken(tc.iat, tc.validity, tc.user, tc.groups, tc.privKey)
	if err != nil {
		log.Fatal("Generating token failed", err)
	}

	err1, user1, groups1 := verifyToken(token, &[]SettingsCredential{{Key: tc.pubKey, UserGlobs: []string{"temper", "temporal"}, GroupGlobs: []string{"temper*"}}})
	if err1 != nil {
		return err1, "", nil
	}
	return nil, user1, groups1
}

func (tc *JWTTestCase) compareIdentity(user1 string, groups1 []string) bool {
	if tc.user != user1 {
		tc.t.Log("Users did not match")
		return false
	}
	if len(groups1) != len(tc.groups) {
		tc.t.Log("Number of granted groups does not match")
		return false
	}
	for i, group := range tc.groups {
		if group != groups1[i] {
			tc.t.Log("Groups did not match")
			return false
		}
	}
	return true
}

func TestHappyCase(t *testing.T) {
	user, groups := "temper", []string{"temper1", "temper2"}
	tc := newJWTTestCase(time.Now().Add(-1*time.Second), 60*time.Second, user, groups, user, "temper*", t)
	err, user1, groups1 := tc.createAndCheckToken()
	if err != nil {
		t.Error("Token could not be validated", err)
	}
	if !tc.compareIdentity(user1, groups1) {
		t.Error("Identity does not match")
	}
}

func TestExpired(t *testing.T) {
	user, groups := "temper", []string{"temper1", "temper2"}
	tc := newJWTTestCase(time.Now().Add(-61*time.Minute), 60*time.Minute, user, groups, user, "temper*", t)
	err, _, _ := tc.createAndCheckToken()
	if err == nil {
		t.Error("Token was accepted even though it was expired")
	}
}

func TestUserGlob(t *testing.T) {
	user, groups := "tempel", []string{"temper1", "temper2"}
	tc := newJWTTestCase(time.Now().Add(-1*time.Second), 60*time.Second, user, groups, "temper*", "temper*", t)
	err, _, _ := tc.createAndCheckToken()
	if err == nil {
		t.Error("Expected user glob to mismatch")
	}
}

func TestGroupGlob(t *testing.T) {
	user, groups := "temper3", []string{"temper1", "temper2"}
	tc := newJWTTestCase(time.Now().Add(-1*time.Second), 60*time.Second, user, groups, "temper*", "tempel*", t)
	err, _, _ := tc.createAndCheckToken()
	if err == nil {
		t.Error("Expected group glob to mismatch")
	}
}

func TestTimeTravel(t *testing.T) {
	user, groups := "temper3", []string{"temper1", "temper2"}
	tc := newJWTTestCase(time.Now().Add(10*time.Second), 60*time.Second, user, groups, "temper*", "temper*", t)
	err, _, _ := tc.createAndCheckToken()
	if err == nil {
		t.Error("Expected time skew to be detected and forbidden")
	}
}

func TestRejectWrongKey(t *testing.T) {
	user, groups := "temper3", []string{"temper1", "temper2"}
	tc := newJWTTestCase(time.Now().Add(-1*time.Second), 60*time.Second, user, groups, "temper*", "temper*", t)
	dat, err := ioutil.ReadFile("./test_resources/RS256-other.test.key")
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(dat)
	if err != nil {
		log.Fatal("Couldn't load private key: ", err)
	}
	tc.privKey = privateKey
	err, _, _ = tc.createAndCheckToken()
	if err == nil {
		t.Error("Expected time skew to be detected and forbidden")
	}
}
