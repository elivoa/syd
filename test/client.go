package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var (
	config = oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},

		RedirectURL: "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8880/auth/authorize",
			TokenURL: "http://localhost:8880/auth/token",
		},

		// RedirectURL: "http://localhost:9094/oauth2",
		// Endpoint: oauth2.Endpoint{
		// 	AuthURL:  "http://localhost:9096/authorize",
		// 	TokenURL: "http://localhost:9096/token",
		// },
	}
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := config.AuthCodeURL("xyz")
		fmt.Println("++++++++++++++++++++++++++++++u: ", u)
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		log.Println("---------1")
		state := r.Form.Get("state")
		if state != "xyz" {
			log.Println("---------2")
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}
		code := r.Form.Get("code")
		if code == "" {
			log.Println("---------3")
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}
		log.Println("---------4", context.Background(), code)
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			log.Println("---------5 - err.Error()", token, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("---------6")
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(*token)
	})

	log.Println("Client is running at 9094 port.")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
