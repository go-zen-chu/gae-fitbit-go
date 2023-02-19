package fitbit2gcal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	da "github.com/go-zen-chu/gae-fitbit-go/domain/auth"
	df2g "github.com/go-zen-chu/gae-fitbit-go/domain/fitbit2gcal"
	"github.com/pkg/errors"
)

const (
	fitbitAPIVersion = "1.2"
)

type fitbitClient struct {
	store  da.Store
	config *df2g.FitbitConfig
}

func NewFitbitClient(store da.Store, config *df2g.FitbitConfig) df2g.FitbitClient {
	return &fitbitClient{
		store:  store,
		config: config,
	}
}

func (fc *fitbitClient) GetSleepData(dateStr string) (*df2g.Sleep, error) {
	token, err := fc.store.FetchFitbitToken()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	tokenSource := fc.config.OauthConfig.TokenSource(ctx, token)
	cli := oauth2.NewClient(ctx, tokenSource)
	defer func() {
		err := fc.updateStoredToken(tokenSource)
		if err != nil {
			log.Errorf("%v\n", err)
		}
	}()
	s, err := requestSleepData(cli, dateStr)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (fc *fitbitClient) GetActivityData(dateStr string) (*df2g.Activity, error) {
	token, err := fc.store.FetchFitbitToken()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	tokenSource := fc.config.OauthConfig.TokenSource(ctx, token)
	cli := oauth2.NewClient(ctx, tokenSource)
	defer func() {
		err := fc.updateStoredToken(tokenSource)
		if err != nil {
			log.Errorf("%v\n", err)
		}
	}()
	a, err := requestActivityData(cli, dateStr)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// make sure to save new token refreshed via oauth2 library
func (fc *fitbitClient) updateStoredToken(tokenSource oauth2.TokenSource) error {
	latestToken, err := tokenSource.Token()
	if err != nil {
		return errors.Wrap(err, "Error getting latest token")
	}
	if err = fc.store.WriteFitbitToken(latestToken); err != nil {
		return errors.Wrap(err, "Error writing latest token")
	}
	return nil
}

func requestSleepData(cli *http.Client, dateStr string) (*df2g.Sleep, error) {
	resBytes, err := requestFitbitData(cli, "sleep", dateStr)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting data")
	}
	var s df2g.Sleep
	err = json.Unmarshal(resBytes, &s)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling sleep data")
	}
	return &s, nil
}

func requestActivityData(cli *http.Client, dateStr string) (*df2g.Activity, error) {
	resBytes, err := requestFitbitData(cli, "activities", dateStr)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting data")
	}
	var a df2g.Activity
	err = json.Unmarshal(resBytes, &a)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling activity data")
	}
	return &a, nil
}

func requestFitbitData(cli *http.Client, dataType, dateStr string) ([]byte, error) {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.fitbit.com"
	u.Path = fmt.Sprintf("/%s/user/-/%s/date/%s.json",
		fitbitAPIVersion,
		dataType,
		dateStr)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}
	res, err := cli.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while requesting data to Fitbit")
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("%v %s\n", res, string(resBytes))
		return nil, errors.Wrap(err, "Error reading response body")
	}
	if res.StatusCode != http.StatusOK {
		log.Errorf("%v %s\n", res, string(resBytes))
		return nil, errors.New("Error response from Fitbit API")
	}
	log.Infof("%v\n", string(resBytes))
	return resBytes, nil
}
