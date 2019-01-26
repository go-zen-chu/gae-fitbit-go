package fitbitauth

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"

	"github.com/pkg/errors"

	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
)

type fileStore struct{}

func NewFileStore() dfba.Store {
	return &fileStore{}
}

func (fs *fileStore) WriteFitbitToken(token *oauth2.Token) error {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return errors.Wrap(err, "Error while marshaling")
	}
	err = ioutil.WriteFile("fitbit_oauth_token.json", tokenBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "Error while writing token to file")
	}
	return nil
}

func (fs *fileStore) FetchFitbitToken() (*oauth2.Token, error) {
	tokenBytes, err := ioutil.ReadFile("fitbit_oauth_token.json")
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
