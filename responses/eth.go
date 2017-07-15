package responses

import (
	"encoding/json"
	"github.com/perthgophers/puddle/messagerouter"
	"io/ioutil"
	"net/http"
	"strconv"
)

// EthTicker will return a string containing the current rate in USD
func EthTicker(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {

	type HTTPResponse struct {
		USD float64 `json:"usd"`
	}

	resp, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD")
	if err != nil {
		w.WriteError("Error: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	var httpResponse HTTPResponse
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteError("Error: " + err.Error())
		return err
	}
	err = json.Unmarshal(b, &httpResponse)
	if err != nil {
		w.WriteError("Error: " + err.Error())
		return err
	}
	w.Write(strconv.FormatFloat(httpResponse.USD, 'f', -1, 64) + "USD")
	return nil
}

func init() {
	Handle("!eth", EthTicker)
}
