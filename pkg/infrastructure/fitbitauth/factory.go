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

func (f *factory) CloudStorageStore(bucketName string) (dfba.Store, error) {
	css, err := NewCloudStorageStore(bucketName)
	if err != nil {
		return nil, err
	}
	return css, nil
}

func (f *factory) FitbitAuthHandler(store dfba.Store, config *oauth2.Config) dfba.FitbitAuthHandler {
	return dfba.NewFitbitAuthHandler(f, store, config)
}

func (f *factory) OAuthClient(config *oauth2.Config) dfba.OAuthClient {
	return NewOAuthClient(config)
}
