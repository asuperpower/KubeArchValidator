package registry

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAcrAccessToken(acrName string) (string, error) {
	// Get the AAD token using Azure's Managed Identity
	aadToken, err := getAadToken()
	if err != nil {
		return "", err
	}

	// Exchange the AAD token for an ACR access token
	acrAccessToken, err := exchangeAadTokenForAcrAccessToken(aadToken, acrName)
	if err != nil {
		return "", err
	}

	return acrAccessToken, nil
}

func getAadToken() (string, error) {
	// TODO: Implement this function to obtain an AAD token using Azure's Managed Identity
	return "", nil
}

func exchangeAadTokenForAcrAccessToken(aadToken string, acrName string) (string, error) {
	url := fmt.Sprintf("https://%s.azurecr.io/oauth2/exchange", acrName)

	data := url.Values{}
	data.Set("grant_type", "access_token")
	data.Set("service", fmt.Sprintf("%s.azurecr.io", acrName))
	data.Set("access_token", aadToken)

	resp, err := http.PostForm(url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Check that the access_token field is present
	if _, ok := result["access_token"]; !ok {
		return "", fmt.Errorf("access_token field missing in response")
	}

	return result["access_token"].(string), nil
}
