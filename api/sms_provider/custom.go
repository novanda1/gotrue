package sms_provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/netlify/gotrue/conf"
)

type CustomProvider struct {
	Config  *conf.CustomProviderConfiguration
	APIPath string
}

type customProviderErrResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (t customProviderErrResponse) Error() string {
	return fmt.Sprintf("%s", t.Message)
}

func NewCustomProvider(config conf.CustomProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	apiPath := config.Url
	return &CustomProvider{
		Config:  &config,
		APIPath: apiPath,
	}, nil
}

func (t *CustomProvider) SendSms(phone string, message string) error {
	body := url.Values{
		"To":     {phone},
		"Body":   {message},
		"Secret": {t.Config.Secret},
	}

	client := &http.Client{Timeout: defaultTimeout}

	r, err := http.NewRequest("POST", t.APIPath, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	if res.StatusCode/100 != 2 {
		resp := &customProviderErrResponse{}
		if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
			return err
		}

		return resp
	}
	defer res.Body.Close()

	return nil
}
