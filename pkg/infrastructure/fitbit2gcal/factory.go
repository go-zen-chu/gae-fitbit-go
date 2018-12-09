package fitbit2gcal

import (
	df2g "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbit2gcal"
	if2g "github.com/go-zen-chu/gae-fitbit-go/pkg/infrastructure/fitbit2gcal"
)

type factory struct{}

func NewFactory() df2g.Factory {
	return &factory{}
}

func (f *factory) Service() df2g.Service {
	return if2g.NewService()
}
