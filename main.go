package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var BEARER_TOKEN string
var CONFIG_PATH string

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ready to Serve Crypto Ticker!")
}

func main() {
	log.SetOutput(os.Stdout)
	BEARER_TOKEN = os.Getenv("BEARER_TOKEN")
	CONFIG_PATH = os.Getenv("CONFIG_PATH")

	if len(CONFIG_PATH) == 0 {
		CONFIG_PATH = "config/"
	}

	if len(BEARER_TOKEN) == 0 {
		fmt.Println("BEARER TOKEN MUST BE SET. Exiting.")
		return
	}

	cryptoClient := NewCryptoClient()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/coins", cryptoClient.coinsHandler)
	fmt.Printf("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
