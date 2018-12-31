//go:generate mockgen -source factory.go -destination mock_factory.go
package fitbit2gcal

// Factory : Creates objects in this domain
type Factory interface {
	Service(gcalConfig *GCalConfig) Service
	FileStore() Store
	FitbitClient(store Store) FitbitClient
	GCalClient(store Store, gcalConfig *GCalConfig) GCalClient
}
