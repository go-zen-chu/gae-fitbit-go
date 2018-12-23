package fitbit2gcal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/pkg/errors"
)

type fitbitClient struct {
	store dfba.Store
}

func NewFitbitClient(store dfba.Store) df2g.FitbitClient {
	return &fitbitClient{
		store: store,
	}
}

func (fc *fitbitClient) GetSleepData(dateStr string) (*df2g.Sleep, error) {
	ft, err := fc.store.FetchFitbitTokens()
	if err != nil {
		return nil, err
	}
	s, err := requestSleepData(ft, dateStr)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (fc *fitbitClient) GetActivityData(dateStr string) (*df2g.Activity, error) {
	ft, err := fc.store.FetchFitbitTokens()
	if err != nil {
		return nil, err
	}
	a, err := requestActivityData(ft, dateStr)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func requestSleepData(ft *dfba.FitbitTokens, dateStr string) (*df2g.Sleep, error) {
	resBytes, err := requestFitbitData(ft, "sleep", dateStr)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting data:")
	}
	var s df2g.Sleep
	err = json.Unmarshal(resBytes, &s)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling sleep data:")
	}
	return &s, nil
}

func requestActivityData(ft *dfba.FitbitTokens, dateStr string) (*df2g.Activity, error) {
	resBytes, err := requestFitbitData(ft, "activities", dateStr)
	if err != nil {
		return nil, errors.Wrap(err, "Error requesting data:")
	}
	var a df2g.Activity
	err = json.Unmarshal(resBytes, &a)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling activity data:")
	}
	return &a, nil
}

func requestFitbitData(ft *dfba.FitbitTokens, dataType, dateStr string) ([]byte, error) {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.fitbit.com"
	u.Path = fmt.Sprintf("/1/user/%s/%s/date/%s.json", ft.UserID, dataType, dateStr)

	req, err := http.NewRequest("Get", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request: ")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", ft.AccessToken))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while requesting data to Fitbit: ")
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading response body: ")
	}
	return resBytes, nil
}
