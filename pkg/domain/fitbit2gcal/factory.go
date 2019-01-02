//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbit2gcal

// Factory : Creates objects in this domain
type Factory interface {
	Service(gcalConfig *GCalConfig) Service
	FileStore() Store
	FitbitClient(store Store) FitbitClient
	GCalClient(store Store, gcalConfig *GCalConfig) GCalClient
}
