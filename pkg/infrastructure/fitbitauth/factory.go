package fitbitauth

import (
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"golang.org/x/oauth2"
)

type factory struct{}

func NewFactory() dfba.Factory {
	return &factory{}
}

func (f *factory) FileStore() (dfba.Store, error) {
	fs := NewFileStore()
	return fs, nil
}

func (f *factory) FitbitAuthHandler(config *oauth2.Config) dfba.FitbitAuthHandler {
	return dfba.NewFitbitAuthHandler(f, config)
}

func (f *factory) OAuthClient(config *oauth2.Config) dfba.OAuthClient {
	return NewOAuthClient(config)
}