package responses

import (
	"encoding/json"
	"github.com/perthgophers/puddle/messagerouter"
	"io/ioutil"
	"net/http"
)

// BitcoinTicker will return a string containing the current rate in USD
func BitcoinTicker(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {

	type HTTPResponse struct {
		Last string `json:"last"`
	}

	resp, err := http.Get("https://www.bitstamp.net/api/ticker/")
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
	w.Write(string(httpResponse.Last) + "USD")
	return nil
}

func init() {
	Handle("!bitcoin", BitcoinTicker)
}
