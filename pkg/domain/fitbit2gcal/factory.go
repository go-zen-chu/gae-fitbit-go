//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbit2gcal

import (
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	dga "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/gcalauth"
)

// Factory : Creates objects in this domain
type Factory interface {
	Service(fitbitConfig *FitbitConfig, gcalConfig *GCalConfig) Service
	FitbitFileStore() dfba.Store
	GCalFileStore() dga.Store
	FitbitClient(store dfba.Store, fitbitConfig *FitbitConfig) FitbitClient
	GCalClient(store dga.Store, gcalConfig *GCalConfig) GCalClient
}
