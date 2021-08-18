package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Config represents all configuration options for the ticker
type Config struct {
	Ticker struct {
		UUID       string   `json:"uuid"`
		VsCurrency string   `json:"vs_currency"`
		TellJokes  bool     `json:"tell_jokes"`
		Crypto     []string `json:"crypto"`
	} `json:"ticker"`
}

var BEARER_TOKEN string

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ready to Serve Crypto Ticker!")
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		loadConfig(r, w)
	case "POST":
		token := r.Header.Get("Authorization")
		fmt.Println(token)
		if token != fmt.Sprintf("Bearer %v", BEARER_TOKEN) {
			http.Error(w, "Incorrect Credentials", http.StatusUnauthorized)
			return
		}
		config := Config{}
		err := json.NewDecoder(r.Body).Decode(&config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json, err := json.Marshal(config)
		err = ioutil.WriteFile(fmt.Sprintf("config/%v.json", config.Ticker.UUID), json, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Configuration successfully saved."))
	}
}

func loadConfig(r *http.Request, w http.ResponseWriter) {
	query := r.URL.Query()
	uuid := query.Get("UUID")

	if len(uuid) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("A UUID must be provided"))
	} else {
		fileName := fmt.Sprintf("config/%v.json", uuid)
		fmt.Println(fileName)
		if _, err := os.Stat(fileName); err == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			http.ServeFile(w, r, fileName)
		} else if os.IsNotExist(err) {
			w.WriteHeader(404)
			w.Write([]byte("Config does not exist for the provided UUID"))
		} else {
			w.WriteHeader(404)
			w.Write([]byte("Schrodinger config. Don't open the box"))
		}
	}
}

func main() {
	log.SetOutput(os.Stdout)
	BEARER_TOKEN = os.Getenv("BEARER_TOKEN")
	if len(BEARER_TOKEN) == 0 {
		fmt.Println("BEARER TOKEN MUST BE SET. Exiting.")
		return
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/config", configHandler)
	fmt.Printf("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
