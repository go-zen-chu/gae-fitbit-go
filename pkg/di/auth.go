package di

import (
	"errors"

	da "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/auth"
	ia "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/auth"

	"golang.org/x/oauth2"
)

func (di DI) InitFitbitOAuthConfig(config *oauth2.Config) {
	key := "FitbitOAuthConfig"
	di[key] = config
}

func (di DI) FitbitOAuthConfig() *oauth2.Config {
	key := "FitbitOAuthConfig"
	if di[key] == nil {
		panic(errors.New(key + " is nil. Need to be initialized"))
	}
	return di[key].(*oauth2.Config)
}

func (di DI) FitbitOAuthClient() da.OAuthClient {
	key := "FitbitOAuthClient"
	if di[key] == nil {
		di[key] = da.NewOAuthClient(di.FitbitOAuthConfig())
	}
	return di[key].(*oauth2.Config)
}

func (di DI) InitGCalOAuthConfig(config *oauth2.Config) {
	key := "GCalOAuthConfig"
	di[key] = config
}

func (di DI) GCalOAuthConfig() *oauth2.Config {
	key := "GCalOAuthConfig"
	if di[key] == nil {
		panic(errors.New(key + " is nil. Need to be initialized"))
	}
	return di[key].(*oauth2.Config)
}

func (di DI) GCalOAuthClient() da.OAuthClient {
	key := "GCalOAuthClient"
	if di[key] == nil {
		di[key] = da.NewOAuthClient(di.GCalOAuthConfig())
	}
	return di[key].(*oauth2.Config)
}

func (di DI) InitAuthFileStore() {
	key := "AuthFileStore"

}

func (di DI) AuthFileStore() da.Store {
	key := "AuthFileStore"
	if di[key] == nil {
		di[key] = ia.NewFileStore()
	}
	return di[key].(da.Store)
}

func (di DI) InitAuthCloudStorageStore(bucketName string) {
	key := "AuthCloudStorageStore"
	st, err := ia.NewCloudStorageStore(bucketName)
	if err != nil {
		panic(errors.Wrap(err, "Error while making cloud storage store"))
	}
	di[key] = st
}

func (di DI) AuthCloudStorageStore() da.Store {
	key := "AuthCloudStorageStore"
	if di[key] == nil {
		di[key] = ia.NewCloudStorageStore()
	}
	return di[key].(da.Store)
}

func (di DI) AuthStore() da.Store {
	if di["AuthFileStore"] != nil {
		return di["AuthFileStore"].(da.Store)
	} else if di["AuthCloudStorageStore"] != nil {
		return di["AuthCloudStorageStore"].(da.Store)
	} else {
		panic(errors.New("Store is not initialized"))
	}
}

func (di DI) AuthService() da.Service {
	key := "AuthService"
	if di[key] == nil {
		di[key] = da.NewService(di.AuthStore(), di.FitbitOAuthClient(), di.GCalOAuthClient())
	}
	return di[key].(da.Service)
}
