package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Config represents all configuration options for the ticker
type Config struct {
	Ticker struct {
		VsCurrency string   `json:"vs_currency"`
		TellJokes  bool     `json:"tell_jokes"`
		Crypto     []string `json:"crypto"`
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ready to Serve Crypto Ticker!")
}

func configHandler(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/config", configHandler)
	fmt.Printf("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
