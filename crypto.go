package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	coingecko "github.com/superoo7/go-gecko/v3"
	"github.com/superoo7/go-gecko/v3/types"
)

type Crypto struct {
	client     *coingecko.Client
	coinsCache *cache.Cache
}

var COIN_KEY = "COIN"

func NewCryptoClient() *Crypto {
	c := Crypto{}

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	c.client = coingecko.NewClient(httpClient)
	c.coinsCache = cache.New(30*time.Minute, 30*time.Minute)
	return (&c)

}

func (c Crypto) coinsHandler(w http.ResponseWriter, r *http.Request) {
	coins, err := c.getCoins()
	if err != nil {
		http.Error(w, "Error getting coins", http.StatusInternalServerError)
	}
	json, err := json.Marshal(struct {
		Coins *types.CoinList `json:"coins"`
	}{coins})

	if err != nil {
		http.Error(w, "Error creatings JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(json))
}

func (c Crypto) getCoins() (*types.CoinList, error) {
	data, found := c.coinsCache.Get(COIN_KEY)

	if found {
		coins := data.(*types.CoinList)
		return coins, nil
	} else {
		coins, err := c.client.CoinsList()
		if err != nil {
			return nil, err
		}
		c.coinsCache.Set(COIN_KEY, coins, cache.NoExpiration)
		return coins, nil
	}
}
