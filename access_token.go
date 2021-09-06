package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

//GrantType is the type of the request, refresh or get a new token
type GrantType string

const (
	//RefreshToken used to refresh the token
	RefreshToken GrantType = "refresh_token"
	//AuthorizationCode used to get a new code
	AuthorizationCode GrantType = "authorization_code"
)

func makeAccesTokenRequest(accessToken *AccessToken, grantType GrantType) {
	app := getStreamlabsAppInfo()

	logMessage("making request for access_code using grant_type " + string(grantType))

	var grantTypeLabel string
	var grantTypeVal string
	if grantType == RefreshToken {
		grantTypeLabel = "refresh_token"
		grantTypeVal = accessToken.RefreshToken

	} else if grantType == AuthorizationCode {
		grantTypeLabel = "code"
		grantTypeVal = appCode
	}

	message := map[string]interface{}{
		"grant_type":    grantType,
		"client_id":     app.ClientID,
		"client_secret": app.ClientSecret,
		"redirect_uri":  app.RedirectURI,
		grantTypeLabel:  grantTypeVal,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		logErrorMessage("error marshaling request body: " + err.Error())
	}

	resp, err := http.Post("https://streamlabs.com/api/v1.0/token",
		"application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		logError(err)
		accessToken = nil
		return
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		logErrorMessage("bad response getting access_token: " + string(body))
		accessToken = nil
		return
	}

	accessToken.TimeBorn = time.Now()

	err = json.NewDecoder(resp.Body).Decode(&accessToken)
	if err != nil {
		logErrorMessage("error decoding response for access token: " + err.Error())
		accessToken = nil
		return
	}
}

func refreshAccessToken(accessToken *AccessToken) {
	makeAccesTokenRequest(accessToken, RefreshToken)
}

func getAccessToken() *AccessToken {
	var accessToken AccessToken
	makeAccesTokenRequest(&accessToken, AuthorizationCode)

	return &accessToken
}
