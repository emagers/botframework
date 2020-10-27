package authclient

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	getPostFunc   func(url string, contentType string, body io.Reader) (*http.Response, error)
	configuration Configuration = Configuration{
		"https://testing.com",
		"ApplicationID",
		"AccessToken",
	}
)

type httpClientMock struct {
	PostFunc func(url string, contentType string, body io.Reader) (*http.Response, error)
}

func (m *httpClientMock) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	return getPostFunc(url, contentType, body)
}

// TestGetAccessTokenReturnsTokenAfterMakingInitialRequest tests that a new token is returned after making a call to authenticate
func TestGetAccessTokenReturnsTokenAfterMakingInitialRequest(t *testing.T) {
	authClient := Init(configuration)

	Client = &httpClientMock{}

	json := `{"expires_in":3600,"access_token":"sample token"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	getPostFunc = func(url string, contentType string, body io.Reader) (*http.Response, error) {
		assert.EqualValues(t, "https://testing.com", url)
		assert.EqualValues(t, "application/x-www-form-urlencoded", contentType)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	token, err := authClient.GetAccessToken()

	assert.EqualValues(t, "sample token", token)
	assert.Nil(t, err)
}

// TestGetAccessTokenReturnsUnexpiredTokenAfterMakingRequestTwice tests that the first token, which is still valid, will be returned when GetAccessToken is called a second time
func TestGetAccessTokenReturnsUnexpiredTokenAfterMakingRequestTwice(t *testing.T) {
	authClient := Init(configuration)

	Client = &httpClientMock{}

	getPostFunc = func(url string, contentType string, body io.Reader) (*http.Response, error) {
		json := `{"expires_in":3600,"access_token":"sample token"}`
		r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	token, err := authClient.GetAccessToken()

	getPostFunc = func(url string, contentType string, body io.Reader) (*http.Response, error) {
		json := `{"expires_in": 3600, "access_token": "second token"}`
		r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	token, err = authClient.GetAccessToken()

	assert.EqualValues(t, "sample token", token)
	assert.Nil(t, err)
}

// TestGetAccessTokenReturnsError tests that when the Azure Bot Framework returns an error, that error is passed back to the calle
func TestGetAccessTokenReturnsError(t *testing.T) {
	authClient := Init(configuration)

	Client = &httpClientMock{}

	msg := "error message"
	getPostFunc = func(url string, contentType string, body io.Reader) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(msg)))

		return &http.Response{
			StatusCode: 400,
			Body:       r,
		}, errors.New(msg)
	}

	token, err := authClient.GetAccessToken()
	assert.Empty(t, token)
	assert.NotNil(t, err)
	assert.EqualValues(t, msg, err.Error())
}
