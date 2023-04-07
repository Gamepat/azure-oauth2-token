package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func RequestAccessTokenInfo(tenantId string, config AuthConfig) (*TokenResponseSuccess, error) {
	if tenantId == "" {
		return nil, fmt.Errorf("tenantId is empty")
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
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(url.QueryEscape(config.ClientID), url.QueryEscape(config.ClientSecret))

	// post request
	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error on request: %v", err)
	}
	defer httpResponse.Body.Close()

	responseBuffer, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http response body: %v", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		// Pretty format error output
		var prettyJSON bytes.Buffer
		json.Indent(&prettyJSON, responseBuffer, "", "    ")
		return nil, fmt.Errorf("response returned with status-code %d and body:\n%s", httpResponse.StatusCode, prettyJSON.String())
	}

	var response TokenResponseSuccess
	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		return nil, fmt.Errorf("error on unmarshal: %v", err)
	}

	return &response, nil
}

func RequestAccessToken(tenantId string, config AuthConfig) (string, error) {

	response, err := RequestAccessTokenInfo(tenantId, config)
	if err != nil {
		return "", err
	}

	if response.AccessToken == "" {
		return "", fmt.Errorf("received access-token is empty")
	}

	return response.AccessToken, nil
}

// ===== Types =====
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	ClientScope  string
}

type TokenResponseSuccess struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
