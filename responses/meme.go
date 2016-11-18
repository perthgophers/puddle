package responses

import (
	"encoding/json"
	"fmt"
	"github.com/perthgophers/puddle/messagerouter"
	"io/ioutil"
	"net/http"
	"net/url"
)

type cookiejar struct {
	jar map[string][]*http.Cookie
}

// Used to set cookies into HTTP lib cookiejar
func (p *cookiejar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	fmt.Printf("The URL is %s\n", u.String())
	fmt.Printf("The cookie being set is : %s\n", cookies)
	p.jar[u.Host] = cookies
}

// Used to get cookies from HTTP lib cookiejar
func (p *cookiejar) Cookies(u *url.URL) []*http.Cookie {
	fmt.Printf("The URL is %s\n", u.String())
	fmt.Printf("The cookie being returned is : %s\n", p.jar[u.Host])
	return p.jar[u.Host]
}

// MemeGenerator will post the latest meme from imgur meme gallery
func MemeGenerator(cr *messagerouter.CommandRequest, w messagerouter.ResponseWriter) error {
	// Create cookiejar so we can authenticate to imgur API
	// http://stackoverflow.com/questions/11361431/authenticated-http-client-requests-from-golang

	// HTTP object to capture imgur API responses
	type HTTPResponse struct {
		Data []struct {
			AccountID    int         `json:"account_id"`
			AccountURL   string      `json:"account_url"`
			CommentCount int         `json:"comment_count"`
			Cover        string      `json:"cover"`
			CoverHeight  int         `json:"cover_height"`
			CoverWidth   int         `json:"cover_width"`
			Datetime     int         `json:"datetime"`
			Description  interface{} `json:"description"`
			Downs        int         `json:"downs"`
			Favorite     bool        `json:"favorite"`
			ID           string      `json:"id"`
			ImagesCount  int         `json:"images_count"`
			InGallery    bool        `json:"in_gallery"`
			IsAd         bool        `json:"is_ad"`
			IsAlbum      bool        `json:"is_album"`
			Layout       string      `json:"layout"`
			Link         string      `json:"link"`
			Nsfw         bool        `json:"nsfw"`
			Points       int         `json:"points"`
			Privacy      string      `json:"privacy"`
			Score        int         `json:"score"`
			Section      string      `json:"section"`
			Title        string      `json:"title"`
			Topic        string      `json:"topic"`
			TopicID      int         `json:"topic_id"`
			Ups          int         `json:"ups"`
			Views        int         `json:"views"`
			Vote         interface{} `json:"vote"`
		} `json:"data"`
	}

	client := &http.Client{}
	cjar := &cookiejar{}
	cjar.jar = make(map[string][]*http.Cookie)
	client.Jar = cjar

	// Authenticate
	req, err := http.NewRequest("GET", "https://api.imgur.com/3/g/memes/", nil)
	req.Header.Set("Authorization", "Client-ID 91101181cd13628")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Unable to get auth token from imgur api", err)
	}

	defer resp.Body.Close()

	var httpResponse HTTPResponse
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read", err)
	}

	err = json.Unmarshal(b, &httpResponse)

	if err != nil {
		fmt.Println("Unable to unmarshall", err)
	}

	// Post meme to slack channel using RTM
	w.Write(string(httpResponse.Data[0].Link))

	return nil
}

func init() {
	Handle("!meme", MemeGenerator)
}
