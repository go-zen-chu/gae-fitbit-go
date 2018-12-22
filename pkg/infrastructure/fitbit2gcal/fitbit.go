package fitbit2gcal

import (
	"fmt"
	"net/http"
	"net/url"

	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/pkg/errors"
)

type fitbitClient struct {
	store dfba.Store
}

func (fc *fitbitClient) GetSleepData() error {
	ft, err := fc.store.FetchFitbitTokens()
	if err != nil {
		return errors.Wrap(err, "Error fetching fitbit token: ")
	}

	return nil
}

func (fc *fitbitClient) GetActivityData() error {
	ft, err := fc.store.FetchFitbitTokens()
	if err != nil {
		return errors.Wrap(err, "Error fetching fitbit token: ")
	}
	return nil
}

func requestSleepData() {

}

func requestActivityData(ft *dfba.FitbitTokens, dateStr string) {
	requestFitbitData(ft, "activities", dateStr)

}

func requestFitbitData(ft *dfba.FitbitTokens, dataType, dateStr string) ([]byte, errors) {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.fitbit.com"
	u.Path = fmt.Sprintf("/1/user/%s/%s/date/%s.json", ft.UserID, dataType, dateStr)

	req, err := http.NewRequest("Get", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating request: ")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", accessToken))
	res, err = client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error while requesting data to Fitbit: ")
	}
	return res, nil
}
