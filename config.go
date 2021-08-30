package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type Register struct {
	UUID string `json:"uuid"`
}

func defaultConfig(uuid string) Config {
	defaults := Config{}
	defaults.Ticker.UUID = uuid
	defaults.Ticker.VsCurrency = "usd"
	defaults.Ticker.TellJokes = true
	defaults.Ticker.Crypto = []string{"bitcoin", "dogecoin", "ethereum"}

	return defaults
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Error(w, "Unsupported Method", http.StatusBadRequest)
		return
	case "POST":
		token := r.Header.Get("Authorization")
		if token != fmt.Sprintf("Bearer %v", BEARER_TOKEN) {
			http.Error(w, "Incorrect Credentials", http.StatusUnauthorized)
			return
		}
		device := Register{}
		err := json.NewDecoder(r.Body).Decode(&device)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fileName := fmt.Sprintf("%v%v.json", CONFIG_PATH, device.UUID)
		if _, err := os.Stat(fileName); err == nil {
			http.Error(w, "Device is already registered", http.StatusOK)
			return
		} else if os.IsNotExist(err) {
			newConfig := defaultConfig(device.UUID)
			json, err := json.Marshal(newConfig)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			err = ioutil.WriteFile(fileName, json, 0644)
			w.Write([]byte("Device Registered"))
			return
		} else {
			http.Error(w, "Something bad happened. Try again.", http.StatusInternalServerError)
		}
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		loadConfig(r, w)
	case "POST":
		token := r.Header.Get("Authorization")
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

		fileName := fmt.Sprintf("%v%v.json", CONFIG_PATH, config.Ticker.UUID)
		if _, err := os.Stat(fileName); err == nil {
			err = ioutil.WriteFile(fileName, json, 0644)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("Configuration successfully saved."))
		} else {
			http.Error(w, "Device is not registered. Must be registered before config can be saved.", http.StatusForbidden)
		}
	}
}

func loadConfig(r *http.Request, w http.ResponseWriter) {
	query := r.URL.Query()
	uuid := query.Get("UUID")

	if len(uuid) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("A UUID must be provided"))
	} else {
		fileName := fmt.Sprintf("%v%v.json", CONFIG_PATH, uuid)
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
