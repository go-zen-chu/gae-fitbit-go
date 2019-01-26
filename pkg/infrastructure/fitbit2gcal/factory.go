package fitbit2gcal

import (
	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	ifba "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbitauth"
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
	iga "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/gcalauth"
)

type factory struct{}

func NewFactory() df2g.Factory {
	return &factory{}
}

func (f *factory) Service(fitbitConfig *df2g.FitbitConfig, gcalConfig *df2g.GCalConfig) df2g.Service {
	fbst := f.FitbitFileStore()
	fbc := NewFitbitClient(fbst, fitbitConfig)
	gst := f.GCalFileStore()
	gc := NewGCalClient(gst, gcalConfig)
	return df2g.NewService(fbc, gc)
}

func (f *factory) FitbitFileStore() dfba.Store {
	return ifba.NewFileStore()
}

func (f *factory) GCalFileStore() dga.Store {
	return iga.NewFileStore()
}

func (f *factory) FitbitClient(store dfba.Store, fitbitConfig *df2g.FitbitConfig) df2g.FitbitClient {
	fc := NewFitbitClient(store, fitbitConfig)
	return fc
}

func (f *factory) GCalClient(store dga.Store, gcalConfig *df2g.GCalConfig) df2g.GCalClient {
	gc := NewGCalClient(store, gcalConfig)
	return gc
}
