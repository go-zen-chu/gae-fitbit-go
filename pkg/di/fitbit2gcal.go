package di

import (
	"errors"

	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	if2g "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbit2gcal"
)

func (di DI) FitbitConfig() *df2g.FitbitConfig {
	key := "FitbitConfig"
	if di[key] == nil {
		di[key] = &df2g.FitbitConfig{
			OauthConfig: di.FitbitOAuthConfig(),
		}
	}
	return di[key].(*df2g.FitbitConfig)
}

func (di DI) FitbitClient() df2g.FitbitClient {
	key := "FitbitClient"
	if di[key] == nil {
		di[key] = if2g.NewFitbitClient(di.AuthStore(), di.FitbitConfig())
	}
	return di[key].(df2g.FitbitClient)
}

func (di DI) InitGCalConfig(sleepCalID, activityCalID string) {
	key := "GCalConfig"
	di[key] = &df2g.GCalConfig{
		SleepCalendarID:    sleepCalID,
		ActivityCalendarID: activityCalID,
		OauthConfig:        di.GCalOAuthConfig(),
	}
}

func (di DI) GCalConfig() *df2g.GCalConfig {
	key := "GCalConfig"
	if di[key] == nil {
		panic(errors.New(key + " is nil. Need to be initialized"))
	}
	return di[key].(*df2g.GCalConfig)
}

func (di DI) GCalClient() df2g.GCalClient {
	key := "GCalClient"
	if di[key] == nil {
		di[key] = if2g.NewGCalClient(di.AuthStore(), di.GCalConfig())
	}
	return di[key].(df2g.GCalClient)
}

func (di DI) Fitbit2GCalService() df2g.Service {
	key := "Fitbit2GCalService"
	if di[key] == nil {
		di[key] = df2g.NewService(di.FitbitClient(), di.GCalClient())
	}
	return di[key].(df2g.Service)
}
