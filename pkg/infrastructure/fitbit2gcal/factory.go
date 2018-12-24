package fitbit2gcal

import (
	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
)

type factory struct{}

func NewFactory() df2g.Factory {
	return &factory{}
}

func (f *factory) Service(gcalConfig *df2g.GCalConfig) df2g.Service {
	st := f.FileStore()
	fbc := NewFitbitClient(st)
	gc := NewGCalClient(st, gcalConfig)
	return df2g.NewService(fbc, gc)
}

func (f *factory) FileStore() df2g.Store {
	return NewFileStore()
}

func (f *factory) FitbitClient(store df2g.Store) df2g.FitbitClient {
	fc := NewFitbitClient(store)
	return fc
}

func (f *factory) GCalClient(store df2g.Store, gcalConfig *df2g.GCalConfig) df2g.GCalClient {
	gc := NewGCalClient(store, gcalConfig)
	return gc
}
