package fitbit2gcal

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/pkg/errors"
)

const (
	fitbitAPIVersion = "1.2"
)

type fitbitClient struct {
	store dfba.Store
	config *df2g.FitbitConfig
}

func NewFitbitClient(store dfba.Store, config *df2g.FitbitConfig) df2g.FitbitClient {
	return &fitbitClient{
		store: store,
		config: config,
	}
}

func (fc *fitbitClient) GetSleepData(dateStr string) (*df2g.Sleep, error) {
	token, err := fc.store.FetchFitbitToken()
	if err != nil {
		return nil, err
	}
	cli := fc.config.OauthConfig.Client(context.Background(), token)
	// make sure to save new token refreshed via oauth2 library
	defer func () {
		if err := fc.store.WriteFitbitToken(token); err != nil {
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
	cli := fc.config.OauthConfig.Client(context.Background(), token)
	// make sure to save new token refreshed via oauth2 library
	defer func () {
		if err := fc.store.WriteFitbitToken(token); err != nil {
			log.Errorf("%v\n", err)
		}
	}()
	a, err := requestActivityData(cli, dateStr)
	if err != nil {
		return nil, err
	}
	return a, nil
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
	// trasport will automatically add header
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.)
	log.Debugf("%v\n", req)

	res, err := cli.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while requesting data to Fitbit")
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Debugf("%v %s\n", res, string(resBytes))
		return nil, errors.Wrap(err, "Error reading response body")
	}
	if res.StatusCode != http.StatusOK {
		log.Debugf("%v %s\n", res, string(resBytes))
		return nil, errors.New("Error response from Fitbit API")
	}
	log.Infof("%v\n", string(resBytes))
	return resBytes, nil
}
