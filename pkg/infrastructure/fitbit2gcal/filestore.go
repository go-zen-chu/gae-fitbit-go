package fitbit2gcal

import (
	"encoding/json"
	"io/ioutil"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type fileStore struct{}

func NewFileStore() df2g.Store {
	return &fileStore{}
}

func (fs *fileStore) FetchFitbitTokens() (*dfba.FitbitTokens, error) {
	tokenBytes, err := ioutil.ReadFile("./fitbit-tokens.json")
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting token from file")
	}
	var ft dfba.FitbitTokens
	err = json.Unmarshal(tokenBytes, &ft)
	if err != nil {
		return nil, errors.Wrap(err, "Error while unmarshaling token from file")
	}
	return &ft, nil
}

func (fs *fileStore) FetchGCalTokens() (*oauth2.Token, error) {
	tokenBytes, err := ioutil.ReadFile("./gcal-token.json")
	if err != nil {
		return nil, errors.Wrap(err, "Error while getting token from file")
	}
	var token oauth2.Token
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return nil, errors.Wrap(err, "Error while unmarshaling token from file")
	}
	return &token, nil
}
