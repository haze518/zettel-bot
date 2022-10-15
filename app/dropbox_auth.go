package zettel_bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const AUTH_URL = "https://www.dropbox.com/oauth2/authorize?client_id=%s&response_type=code"
const TOKEN_URL = "https://api.dropboxapi.com/oauth2/token"

type DropboxAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int16  `json:"expires_in"`
	Scope       string `json:"scope"`
	Uid         string `json:"uid"`
	AccountId   string `json:"account_id"`
}

type DropboxAuth struct {
	appKey, clientSecret string
}

func New() *DropboxAuth {
	appKey, ok := os.LookupEnv("APP_KEY")
	if !ok {
		log.Fatal("There is no APP_KEY env var")
	}
	clientSecret, ok := os.LookupEnv("APP_SECRET")
	if !ok {
		log.Fatal("There is no CLIENT_SECRET env var")
	}
	log.Printf("Client credentilas: \nClient Key: %s\nClient Secret: %s\n", appKey, clientSecret)
	return &DropboxAuth{appKey: appKey, clientSecret: clientSecret}
}

func (d *DropboxAuth) getAuthorizationURLMessage() string {
	authUrl := fmt.Sprintf(AUTH_URL, d.appKey)
	return fmt.Sprintf(
		"%s\n%s\n%s",
		"Go to the following URL and allow access",
		"Please send acces code via /access_code {code}",
		authUrl,
	)
}

func (d *DropboxAuth) getAccessToken(authCode string) (*DropboxAuthResponse, error) {
	result := new(DropboxAuthResponse)
	data := url.Values{}
	data.Set("code", authCode)
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", d.appKey)
	data.Set("client_secret", d.clientSecret)
	encodedData := data.Encode()
	req, err := http.NewRequest("POST", TOKEN_URL, strings.NewReader(encodedData))
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Add("Accept", "*/*")
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(body, result)
	return result, nil
}
