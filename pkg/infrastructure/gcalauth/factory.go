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

func (f *factory) CloudStorageStore(bucketName string) (dga.Store, error) {
	css, err := NewCloudStorageStore(bucketName)
	if err != nil {
		return nil, err
	}
	return css, nil
}

func (f *factory) GCalAuthHandler(store dga.Store, oauthConfig *oauth2.Config) dga.GCalAuthHandler {
	return dga.NewGCalAuthHandler(f, store, oauthConfig)
}

func (f *factory) NewOAuthClient(oauthConfig *oauth2.Config) dga.OAuthClient {
	return NewOAuthClient(oauthConfig)
}

func (f *factory) OAuthClient(config *oauth2.Config) dga.OAuthClient {
	return NewOAuthClient(config)
}