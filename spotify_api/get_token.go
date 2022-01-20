package spotify_api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func stringToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (c *Client) getToken() (string, error) {
	if c.token.expires != nil && c.token.expires.After(time.Now()) {
		return c.token.value, nil
	}
	err := c.UpdateToken()
	if err != nil {
		return "", err
	}

	return c.token.value, nil
}

func (c *Client) UpdateToken() error {
	const tokenUrl = "https://accounts.spotify.com/api/token"

	log.Println("Getting new token...")
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Basic "+stringToBase64(fmt.Sprintf("%v:%v", c.clientId, c.clientSecret)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	hc := http.Client{}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	token := &struct {
		Token     string `json:"access_token"`
		ExpiresIn int    `json:"expires_in"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(token)
	if err != nil {
		return err
	}

	expireDate, err := time.ParseDuration(fmt.Sprint(token.ExpiresIn) + "s")
	if err != nil {
		return err
	}

	newExpiresDate := time.Now().Add(expireDate)

	c.token.expires = &newExpiresDate
	c.token.value = token.Token

	log.Println("Just got new token...")
	return nil
}
