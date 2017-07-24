package main

import (
	"encoding/json"
	"log"
	"net/http"
	authentication "k8s.io/kubernetes/pkg/apis/authentication/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func authenticationResponse(w http.ResponseWriter, trs *authentication.TokenReviewStatus) {
	tr := authentication.TokenReview{
		TypeMeta: metav1.TypeMeta{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind: "TokenReview"},
		Status: *trs,
		}
	json.NewEncoder(w).Encode(tr)
}

func runServer(config *Settings){
	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var tr authentication.TokenReview
		err := decoder.Decode(&tr)
		if err != nil {
			log.Println("[Error]", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			authenticationResponse(w,&authentication.TokenReviewStatus{Authenticated:false})
			return
		}
		/*// Check User
		if err != nil {
			//Reject case
			log.Println("[Error]", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			authenticationResponse(w,&authentication.TokenReviewStatus{Authenticated:false})
			return
		}*/

		//Accept case
		log.Println(tr.Spec.Token)
		err, user, groups := verifyToken(tr.Spec.Token,config.Credentials)
		if err != nil {
			log.Println("[Error]", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			// Not responding with error because that could leak information.
			authenticationResponse(w,&authentication.TokenReviewStatus{Authenticated:false,
				Error: "Access not granted"})
		}
		log.Printf("[Success] login as %s", user)
		w.WriteHeader(http.StatusOK)
		trs := authentication.TokenReviewStatus{
			Authenticated: true,
			User: authentication.UserInfo{
				Username: user,
				UID:      user,
				Groups:	  groups,
			},
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"apiVersion": "authentication.k8s.io/v1beta1",
			"kind":       "TokenReview",
			"status":     trs,
		})
	})
	log.Fatal(http.ListenAndServe(":3000", nil))

}
