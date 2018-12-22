package fitbitauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// FitbitAuthHandler : Manage any authorization process against Fitbit API
type FitbitAuthHandler interface {
	Redirect2Fitbit(w http.ResponseWriter, r *http.Request)
	HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request)
	GetFitbitTokens(authCode string, fbClientSecret string) (*FitbitTokens, error)
}

type fitbitAuthHandler struct {
	fitbitAuthParams  *FitbitAuthParams
	fitbitTokenParams *FitbitTokenParams
	factory           Factory
}

func NewFitbitAuthHandler(fac *FitbitAuthParams, ftp *FitbitTokenParams, fbaf Factory) FitbitAuthHandler {
	return &fitbitAuthHandler{
		fitbitAuthParams:  fac,
		fitbitTokenParams: ftp,
		factory:           fbaf,
	}
}

// Redirect to Fitbit's oauth url
func (fah *fitbitAuthHandler) Redirect2Fitbit(w http.ResponseWriter, r *http.Request) {
	redirectURL := CreateFitbitAuthURL(fah.fitbitAuthParams)
	http.Redirect(w, r, redirectURL.String(), http.StatusSeeOther)
}

// CreateFitbitAuthURL : Generate Fitbit's oauth authorization flow url
func CreateFitbitAuthURL(fac *FitbitAuthParams) *url.URL {
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

// HandleFitbitAuthCode : Will recieve auth code from Fitbit, store it to Store
func (fah *fitbitAuthHandler) HandleFitbitAuthCode(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["code"]

	if !ok || len(keys[0]) < 1 {
		log.Errorf("Could not get auth code from request")
		return
	}

	st, err := fah.factory.FileStore()
	if err != nil {
		log.Errorf("Error while getting store : %v", err)
		return
	}

	err = st.WriteAuthCode(keys[0])
	if err != nil {
		log.Errorf("Error while writing auth code : %v", err)
		return
	}
}

func (fah *fitbitAuthHandler) GetFitbitTokens(authCode string, fbClientSecret string) (*FitbitTokens, error) {
	base64secret := base64.StdEncoding.EncodeToString([]byte(fbClientSecret))
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "www.fitbit.com"
	u.Path = "/oauth2/token"

	v := url.Values{}
	v.Set("code", authCode)
	v.Add("grant_type", fah.fitbitTokenParams.GrantType)
	v.Add("client_id", fah.fitbitTokenParams.ClientID)
	v.Add("redirect_uri", fah.fitbitTokenParams.RedirectURI)

	rd := strings.NewReader(v.Encode())
	req, err := http.NewRequest("POST", u.String(), rd)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64secret))
	req.Header.Set("content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fj := &FitbitTokens{}
	if err = json.NewDecoder(res.Body).Decode(fj); err != nil {
		return nil, err
	}
	return fj, nil
}
