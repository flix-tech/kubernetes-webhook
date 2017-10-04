package main

import (
	"encoding/json"
	"log"
	"net/http"
	authentication "k8s.io/kubernetes/pkg/apis/authentication/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func authenticationResponse(w http.ResponseWriter, trs *authentication.TokenReviewStatus) {
	tr := authentication.TokenReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "authentication.k8s.io/v1beta1",
			Kind:       "TokenReview"},
		Status: *trs,
	}
	json.NewEncoder(w).Encode(tr)
}

func authenticate(w http.ResponseWriter, r *http.Request, config *Settings) {
	decoder := json.NewDecoder(r.Body)
	var tr authentication.TokenReview
	err := decoder.Decode(&tr)
	if err != nil {
		log.Println("[Error]", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		authenticationResponse(w, &authentication.TokenReviewStatus{Authenticated: false, Error: "Bad Request"})
		return
	}

	// Main verificatoin step
	err, user, groups := verifyToken(tr.Spec.Token, config.Credentials)
	if err != nil {
		log.Println("[Error]", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		// Not responding with error because that could leak information.
		authenticationResponse(w, &authentication.TokenReviewStatus{Authenticated: false,
			Error: "Access not granted"})
		return
	}
	log.Printf("[Success] login as %s", user)
	w.WriteHeader(http.StatusOK)
	trs := authentication.TokenReviewStatus{
		Authenticated: true,
		User: authentication.UserInfo{
			Username: user,
			UID:      user,
			Groups:   groups,
		},
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"apiVersion": "authentication.k8s.io/v1beta1",
		"kind":       "TokenReview",
		"status":     trs,
	})
}

func assertFileExists(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("File '%v' does not exist.", path)
		os.Exit(1)
	}
}

func runServer(config *Settings) *http.Server {
	srv := &http.Server{Addr: config.ListenAddress}
	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		authenticate(w, r, config)
	})
	if config.SSL {
		assertFileExists(config.SSLKeyPath)
		assertFileExists(config.SSLCrtPath)
		log.Printf("Starting HTTPS server on address: %v", srv.Addr)
		go func() { log.Fatal(srv.ListenAndServeTLS(config.SSLCrtPath, config.SSLKeyPath)) }()
	} else {
		log.Printf("Starting HTTP server on address: %v", srv.Addr)
		go func() { log.Fatal(srv.ListenAndServe()) }()
	}
	return srv
}
