// Package infobip provides a client library
// for interacting with the Infobip API.
// https://www.infobip.com/docs/api
package infobip

import (
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/infobip-community/infobip-api-go-sdk/v3/internal"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/account"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/email"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/mms"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/numbers"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/rcs"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/sms"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/webrtc"
	"github.com/infobip-community/infobip-api-go-sdk/v3/pkg/infobip/whatsapp"
)

// Client is the entrypoint to all Infobip channels and platform.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient http.Client
	WhatsApp   whatsapp.WhatsApp
	MMS        mms.MMS
	Email      email.Email
	SMS        sms.SMS
	WebRTC     webrtc.WebRTC
	RCS        rcs.RCS
	Numbers    numbers.Numbers
	Account    account.Account
}

// NewClientFromEnv returns a client object using the credentials from the environment.
// If a client is not provided using options, a default one is created.
func NewClientFromEnv(options ...func(*Client)) (Client, error) {
	if os.Getenv("IB_BASE_URL") == "" {
		return Client{}, errors.New("IB_BASE_URL environment variable is not set")
	}
	if os.Getenv("IB_API_KEY") == "" {
		return Client{}, errors.New("IB_API_KEY environment variable is not set")
	}

	return NewClient(os.Getenv("IB_BASE_URL"), os.Getenv("IB_API_KEY"), options...)
}

// NewClient returns a client object using the provided baseURL and apiKey.
// If a client is not provided using options, a default one is created.
func NewClient(baseURL string, apiKey string, options ...func(*Client)) (Client, error) {
	baseURL, err := validateURL(baseURL)
	if err != nil {
		return Client{}, err
	}
	c := Client{baseURL: baseURL, apiKey: apiKey, httpClient: http.Client{}}

	for _, opt := range options {
		opt(&c)
	}

	c.WhatsApp = &whatsapp.Channel{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}
	c.MMS = &mms.Channel{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}

	c.Email = &email.Channel{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}

	c.SMS = &sms.Channel{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}

	c.WebRTC = &webrtc.Channel{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}

	c.RCS = &rcs.Channel{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}

	c.Numbers = &numbers.Platform{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}

	c.Account = &account.Platform{
		ReqHandler: internal.HTTPHandler{APIKey: apiKey, BaseURL: baseURL, HTTPClient: c.httpClient},
	}
	return c, nil
}

func validateURL(baseURL string) (string, error) {
	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		baseURL = "https://" + baseURL
		_, err = url.ParseRequestURI(baseURL)
	}

	return baseURL, err
}

func WithHTTPClient(httpClient http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
