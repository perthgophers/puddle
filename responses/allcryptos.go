package responses

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/perthgophers/puddle/messagerouter"
)

// Tickers contain a slice of Ticker
type Tickers []*Ticker

// Ticker contains all the data for a single crypto
type Ticker struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Two4HVolumeUsd   string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	PercentChange1H  string `json:"percent_change_1h"`
	PercentChange24H string `json:"percent_change_24h"`
	PercentChange7D  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

// Get gets the ticker from a slice of tickers
func (ts Tickers) Get(name string) *Ticker {
	for _, v := range ts {
		if v.ID == name {
			return v
		}
	}
	return nil
}

// AllCryptos wraps allCryptos with results
func AllCryptos(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
	if err != nil {
		return err
	}
	result := &Tickers{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return err
	}

	w.Write(formatPrice("stratis", result))
	w.Write(formatPrice("ethereum", result))
	w.Write(formatPrice("bitcoin", result))
	w.Write(formatPrice("bitcoin-cash", result))
	w.Write(formatPrice("neo", result))

	return nil
}

func formatPrice(id string, t *Tickers) string {
	ticker := t.Get(id)
	if ticker == nil {
		log.Println("Could not find ticker")
		return ""
	}

	price, err := strconv.ParseFloat(ticker.PriceUsd, 64)
	if err != nil {
		log.Println(err)
		return ""
	}
	return fmt.Sprintf("%s - %.2f USD", ticker.Name, price)

}

func init() {
	Handle("!all", AllCryptos)
}
