package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func RequestAccessToken(tenantId string, config AuthConfig) (string, error) {

	if tenantId == "" {
		return "", errors.New("tenantId is empty")
	}

	authUrl := "https://login.microsoftonline.com/" + tenantId + "/oauth2/v2.0/token"

	bodyData := url.Values{}
	bodyData.Set("client_id", config.ClientID)
	bodyData.Set("client_secret", config.ClientSecret)
	bodyData.Set("grant_type", "client_credentials")
	bodyData.Set("scope", config.ClientScope)
	body := bodyData.Encode()
	bodyBuffer := bytes.NewBuffer([]byte(body))

	// Create request
	request, err := http.NewRequest(http.MethodPost, authUrl, bodyBuffer)
	if err != nil {
		return "", errors.New("error creating request: " + err.Error())
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(url.QueryEscape(config.ClientID), url.QueryEscape(config.ClientSecret))

	// post request
	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return "", errors.New("error on request: " + err.Error())
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return "", errors.New("response returned with status-code " + fmt.Sprintf("%d", httpResponse.StatusCode))
	}

	responseBuffer, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return "", errors.New("error reading http response body: " + err.Error())
	}
	responseJson := string(responseBuffer)

	// TODO: Check & handle if error-response - https://learn.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-auth-code-flow#error-response-1
	var response tokenResponseSuccess
	err = json.Unmarshal([]byte(responseJson), &response)
	if err != nil {
		return "", errors.New("error on unmarshal: " + err.Error())
	}

	if response.AccessToken == "" {
		return "", errors.New("received access-token is empty")
	}

	// TODO: Maybe return full response for additional info
	return response.AccessToken, nil
}

// ===== Types =====
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	ClientScope  string
}

type tokenResponseSuccess struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
