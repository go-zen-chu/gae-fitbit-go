package fitbitauth

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
)

type fileStore struct{}

func NewFileStore() dfba.Store {
	return &fileStore{}
}

func (fs *fileStore) WriteAuthCode(authCode string) error {
	authCodeBytes := []byte(authCode)
	err := ioutil.WriteFile("fitbit-auth-code.txt", authCodeBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "Error while writing auth code to file")
	}
	return nil
}

func (fs *fileStore) FetchAuthCode() (string, error) {
	authCodeBytes, err := ioutil.ReadFile("./fitbit-auth-code.txt")
	if err != nil {
		return "", errors.Wrap(err, "Error while getting auth code from file")
	}
	return string(authCodeBytes), nil
}

func (fs *fileStore) WriteFitbitTokens(ft *dfba.FitbitTokens) error {
	tokenBytes, err := json.Marshal(ft)
	if err != nil {
		return errors.Wrap(err, "Error while marshaling")
	}
	err = ioutil.WriteFile("fitbit-tokens.json", tokenBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "Error while writing token to file")
	}
	return nil
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
