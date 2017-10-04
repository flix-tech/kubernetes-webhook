package main

import (
	"flag"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
	"log"
	"fmt"
	"bufio"
	"os"
	"io/ioutil"
	"bytes"
	"github.com/flix-tech/kubernetes-webhook/jwt-generator/generator"
)
type arrayFlags []string
var groups arrayFlags

func (i *arrayFlags) String() string {
	var buffer bytes.Buffer
	for _, group := range *i {
		buffer.Write([]byte(group))
		buffer.Write([]byte{','})
	}
	return buffer.String()
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	validity := flag.Duration("validity", 3600 * time.Second, "Duration in seconds in which the token should be valid")
	user := flag.String("user", "", "Username")
	flag.Var(&groups, "group", "Group name (flag can be used multiple times)")
	flag.Usage = func() {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "The secret key will be read from STDIN")
	}
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	privateKeyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}

	if len(*user) == 0{
		log.Fatal("Please define a user name")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}
	tokenString, err := generator.GenerateToken(time.Now(),*validity,*user,groups,privateKey)

	if err != nil {
		log.Fatal(err)
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.

	fmt.Print(tokenString)

}