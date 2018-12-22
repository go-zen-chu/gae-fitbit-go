package fitbitauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// FitbitAuthHandler : Manage any authorization process against Fitbit API
type FitbitAuthHandler interface {
	Redirect2Fitbit(w http.ResponseWriter, r *http.Request)
	HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request)
}

type fitbitAuthHandler struct {
	factory           Factory
	fitbitAuthParams  *FitbitAuthParams
	fitbitTokenParams *FitbitTokenParams
}

func NewFitbitAuthHandler(fbaf Factory, fac *FitbitAuthParams, ftp *FitbitTokenParams) FitbitAuthHandler {
	return &fitbitAuthHandler{
		factory:           fbaf,
		fitbitAuthParams:  fac,
		fitbitTokenParams: ftp,
	}
}

// Redirect to Fitbit's oauth url
func (fah *fitbitAuthHandler) Redirect2Fitbit(w http.ResponseWriter, r *http.Request) {
	redirectURL := createFitbitAuthURL(fah.fitbitAuthParams)
	http.Redirect(w, r, redirectURL.String(), http.StatusSeeOther)
}

// createFitbitAuthURL : Generate Fitbit's oauth authorization flow url
func createFitbitAuthURL(fac *FitbitAuthParams) *url.URL {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "www.fitbit.com"
	u.Path = "/oauth2/authorize"
	q := u.Query()
	q.Set("client_id", fac.ClientID)
	q.Set("redirect_uri", fac.RedirectURI)
	q.Set("scope", fac.Scope)
	q.Set("expires_in", fac.Expires)
	q.Set("response_type", fac.ResponseType)
	u.RawQuery = q.Encode()
	log.Debugf("create fitbit auth url : %s", u.String())
	return u
}

// HandleFitbitAuthCode : Will recieve auth code from Fitbit, store it
func (fah *fitbitAuthHandler) HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		log.Errorf("Could not get auth code from request")
		return
	}
	// auth code is one time, no need to save it
	code := keys[0]
	log.Debugf("auth code :%s", code)
	fst, err := fah.factory.FileStore()
	if err != nil {
		log.Errorf("Error while getting store : %v", err)
		return
	}

	ft, err := requestFitbitTokens(fah.fitbitTokenParams, code)
	if err != nil {
		log.Errorf("Error while getting token : %v", err)
		return
	}

	err = fst.WriteFitbitTokens(ft)
	if err != nil {
		log.Errorf("Error while storing token : %v", err)
		return
	}
	log.Info("Success storing fitbit tokens")
	fmt.Fprintf(w, "OK")
}

func requestFitbitTokens(fbtp *FitbitTokenParams, authCode string) (*FitbitTokens, error) {
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

	var ft FitbitTokens
	if err := json.NewDecoder(res.Body).Decode(ft); err != nil {
		return nil, errors.Wrap(err, "Error while decoding token: ")
	}
	return &ft, nil
}
