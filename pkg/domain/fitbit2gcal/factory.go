package fitbit2gcal

// Factory : Creates objects in this domain
type Factory interface {
	Service() Service
	FileStore() Store
	FitbitClient(store Store) FitbitClient
	GCalClient() GCalClient
}
