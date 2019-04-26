package fitbit2gcal

import (
	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
)

type factory struct{}

func NewFactory() df2g.Factory {
	return &factory{}
}

func (f *factory) Service(fitbitConfig *df2g.FitbitConfig, gcalConfig *df2g.GCalConfig, fbst dfba.Store, gst dga.Store) df2g.Service {
	fbc := NewFitbitClient(fbst, fitbitConfig)
	gc := NewGCalClient(gst, gcalConfig)
	return df2g.NewService(fbc, gc)
}

func (f *factory) FitbitClient(store dfba.Store, fitbitConfig *df2g.FitbitConfig) df2g.FitbitClient {
	fc := NewFitbitClient(store, fitbitConfig)
	return fc
}

func (f *factory) GCalClient(store dga.Store, gcalConfig *df2g.GCalConfig) df2g.GCalClient {
	gc := NewGCalClient(store, gcalConfig)
	return gc
}
