package authclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Configuration contains the properties needed to request a new access token from the Azure Bot Framework
type Configuration struct {
	AppID       string
	AccessToken string
}

type accessToken struct {
	token     string
	expiresAt int64
}

type accessTokenResponse struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

// HTTPClient is a wrapper interface for the http.Post method
type HTTPClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

var (
	// Client is a HTTPClient wrapper
	Client HTTPClient
)

// AuthenticationClient manages the authentication session for the Azure Bot Framework
type AuthenticationClient struct {
	configuration  Configuration
	frameworkToken accessToken
}

// Init creates a new instance of the BotFrameworkConnector
func Init(configuration Configuration) *AuthenticationClient {
	Client = &http.Client{}
	return &AuthenticationClient{
		configuration,
		accessToken{},
	}
}

// GetAccessToken returns a valid access token
func (connect *AuthenticationClient) GetAccessToken() (string, error) {
	if connect.frameworkToken.expiresAt < time.Now().Unix() {
		now := time.Now().Unix()
		reqBody := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&scope=https%%3A%%2F%%2Fapi.botframework.com%%2F.default", connect.configuration.AppID, connect.configuration.AccessToken)

		resp, err := Client.Post("https://login.microsoftonline.com/botframework.com/oauth2/v2.0/token", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(reqBody)))
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		response := accessTokenResponse{}

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&response); err != nil {
			return "", err
		}
		defer resp.Body.Close()

		expiresAt := now + response.ExpiresIn
		connect.frameworkToken = accessToken{
			response.AccessToken,
			expiresAt,
		}
	}

	return connect.frameworkToken.token, nil
}
