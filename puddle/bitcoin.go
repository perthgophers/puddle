package puddle

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"
)

// BitcoinTicker will return a string containing the current rate in USD
func BitcoinTicker(username string, msgText string, msg slack.Msg) error {

	type HTTPResponse struct {
		Last string `json:"last"`
	}

	resp, err := http.Get("https://www.bitstamp.net/api/ticker/")
	if err != nil {
		ErrorMessage("Error: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	var httpResponse HTTPResponse
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorMessage("Error: " + err.Error())
		return err
	}
	err = json.Unmarshal(b, &httpResponse)
	if err != nil {
		ErrorMessage("Error: " + err.Error())
		return err
	}
	SendMessage(string(httpResponse.Last) + "USD")
	return nil
}
