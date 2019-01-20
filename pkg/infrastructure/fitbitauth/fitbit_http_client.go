package fitbitauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type fitbitHTTPClient struct{}

func NewFitbitHTTPClient() dfba.FitbitHTTPClient {
	return &fitbitHTTPClient{}
}

func (fhc *fitbitHTTPClient) GetFitbitToken(fbtp *dfba.FitbitTokenParams, authCode string) (*dfba.FitbitTokens, error) {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.fitbit.com"
	u.Path = "/oauth2/token"

	v := url.Values{}
	v.Set("code", authCode)
	v.Add("grant_type", fbtp.GrantType)
	v.Add("client_id", fbtp.ClientID)
	v.Add("redirect_uri", fbtp.RedirectURI)

	rd := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", u.String(), rd)
	if err != nil {
		return nil, err
	}
	// [Auth header](https://dev.fitbit.com/build/reference/web-api/oauth2/#authorization-header)
	clientAndSecret := fmt.Sprintf("%s:%s", fbtp.ClientID, fbtp.ClientSecret)
	base64secret := base64.StdEncoding.EncodeToString([]byte(clientAndSecret))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64secret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Debugf("Sending request to Fitbit API: %v", req)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	log.Infof("Status Code from Fitbit API: %d", res.StatusCode)

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading body")
	}
	// error response from Fitbit API
	if res.StatusCode != http.StatusOK {
		log.Debugf("%d %s", res.StatusCode, string(bodyBytes))
		msg := fmt.Sprintf("Error response from Fitbit API: %d", res.StatusCode)
		return nil, errors.New(msg)
	}

	var ft dfba.FitbitTokens
	if err := json.Unmarshal(bodyBytes, &ft); err != nil {
		log.Debugf("%d %s\n", res.StatusCode, string(bodyBytes))
		return nil, errors.Wrap(err, "Error while decoding token")
	}
	return &ft, nil
}
