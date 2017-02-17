package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/perthgophers/puddle/messagerouter"
)

type Definition struct {
	Metadata struct {
		Provider string `json:"provider"`
	} `json:"metadata"`
	Results []struct {
		ID             string `json:"id"`
		Language       string `json:"language"`
		LexicalEntries []struct {
			Entries []struct {
				Etymologies         []string `json:"etymologies"`
				GrammaticalFeatures []struct {
					Text string `json:"text"`
					Type string `json:"type"`
				} `json:"grammaticalFeatures"`
				HomographNumber string `json:"homographNumber"`
				Senses          []struct {
					Definitions []string `json:"definitions"`
					Examples    []struct {
						Text string `json:"text"`
					} `json:"examples,omitempty"`
					ID      string   `json:"id"`
					Domains []string `json:"domains,omitempty"`
				} `json:"senses"`
			} `json:"entries"`
			Language        string `json:"language"`
			LexicalCategory string `json:"lexicalCategory"`
			Pronunciations  []struct {
				AudioFile        string   `json:"audioFile"`
				Dialects         []string `json:"dialects"`
				PhoneticNotation string   `json:"phoneticNotation"`
				PhoneticSpelling string   `json:"phoneticSpelling"`
			} `json:"pronunciations"`
			Text string `json:"text"`
		} `json:"lexicalEntries"`
		Type string `json:"type"`
		Word string `json:"word"`
	} `json:"results"`
}

// Define fetches the dictionary definition of a word
func Define(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	messageArray := strings.Split(cr.Text, " ")
	if len(messageArray) != 2 {
		w.Write("Only one word please")
		return nil
	}
	query, err := url.Parse(fmt.Sprintf("https://od-api.oxforddictionaries.com:443/api/v1/entries/en/%s", messageArray[1]))
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", query.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("App_id", "9b9556f7")
	req.Header.Set("App_key", "91f8596e9bb7945c57fca8397d428ab8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	definition := &Definition{}
	json.NewDecoder(resp.Body).Decode(definition)
	if len(definition.Results) == 0 {
		w.Write("No definitions found")
		return nil
	}
	w.Write(definition.Results[0].LexicalEntries[0].Entries[0].Senses[0].Definitions[0])
	return nil
}
func init() {
	Handle("!define", Define)
}
