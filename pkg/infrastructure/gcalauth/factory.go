package gcalauth

import (
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
)

type factory struct{}

func NewFactory() dga.Factory {
	return &factory{}
}

func (f *factory) FileStore() (dga.Store, error) {
	fs := NewFileStore()
	return fs, nil
}

// func (f *factory) FitbitAuthHandler(fap *dfba.FitbitAuthParams, ftp *dfba.FitbitTokenParams) dfba.FitbitAuthHandler {
// 	return dfba.NewFitbitAuthHandler(f, fap, ftp)
// }
