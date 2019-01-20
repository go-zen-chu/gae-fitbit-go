package fitbitauth

import (
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
)

type factory struct{}

func NewFactory() dfba.Factory {
	return &factory{}
}

func (f *factory) FileStore() (dfba.Store, error) {
	fs := NewFileStore()
	return fs, nil
}

func (f *factory) FitbitAuthHandler(fap *dfba.FitbitAuthParams, ftp *dfba.FitbitTokenParams) dfba.FitbitAuthHandler {
	return dfba.NewFitbitAuthHandler(f, fap, ftp)
}

func (f *factory) FitbitHTTPClient() dfba.FitbitHTTPClient {
	return NewFitbitHTTPClient()
}
