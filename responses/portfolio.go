package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/perthgophers/puddle/messagerouter"
	"github.com/pkg/errors"
)

// Portfolio will return a string containing the current rate in USD
func Portfolio(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	if cr.Username != "nii236" {
		return errors.New("invalid username " + cr.Username)
	}
	ethAmt := os.Getenv("ETH")
	btcAmt := os.Getenv("BTC")
	bchAmt := os.Getenv("BCH")

	if ethAmt == "" || btcAmt == "" || bchAmt == "" {
		return errors.New("portfolio not provided in environment variables")
	}
	type HTTPResponse struct {
		Last string `json:"last"`
	}

	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/")
	if err != nil {
		return err
	}
	result := &Tickers{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return err
	}

	ETHticker := result.Symbol("ETH")
	BTCticker := result.Symbol("BTC")
	BCHticker := result.Symbol("BCH")

	if ETHticker == nil || BTCticker == nil || BCHticker == nil {
		return errors.New("could not find tickers from API response")
	}

	ethTotal, err := calcValue(ethAmt, ETHticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}
	btcTotal, err := calcValue(btcAmt, BTCticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}
	bchTotal, err := calcValue(bchAmt, BCHticker.PriceUsd)
	if err != nil {
		return errors.Wrap(err, "could not calculate total value")
	}

	total := ethTotal + btcTotal + bchTotal
	aud, err := usdToAud(total)
	if err != nil {
		return errors.Wrap(err, "could not convert to AUD")
	}
	w.Write(fmt.Sprintf("Your crypto net worth is: %.2f USD (%.2f AUD)", total, aud))

	return nil
}

func calcValue(ownAmt, tickerAmt string) (float64, error) {
	own, err := strconv.ParseFloat(ownAmt, 64)
	if err != nil {
		return 0, err
	}
	tick, err := strconv.ParseFloat(tickerAmt, 64)
	if err != nil {
		return 0, err
	}
	return own * tick, err
}

type forexTicker struct {
	Base  string `json:"base"`
	Date  string `json:"date"`
	Rates struct {
		AUD float64 `json:"AUD"`
	} `json:"rates"`
}

func usdToAud(usd float64) (float64, error) {

	resp, err := http.Get("https://api.fixer.io/latest?symbols=USD,AUD")
	if err != nil {
		return 0, err
	}
	result := &forexTicker{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return 0, err
	}

	converted := result.Rates.AUD * usd
	return converted, nil
}

func init() {
	Handle("!portfolio", Portfolio)
}
