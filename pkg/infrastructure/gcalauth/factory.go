package gcalauth

import (
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
	"golang.org/x/oauth2"
)

type factory struct{}

func NewFactory() dga.Factory {
	return &factory{}
}

func (f *factory) FileStore() (dga.Store, error) {
	fs := NewFileStore()
	return fs, nil
}

func (f *factory) GCalAuthHandler(oauthConfig *oauth2.Config) dga.GCalAuthHandler {
	return dga.NewGCalAuthHandler(f, oauthConfig)
}
