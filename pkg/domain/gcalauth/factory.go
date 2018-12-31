//go:generate mockgen -source factory.go -destination mock_factory.go
package gcalauth

// Factory : Creates objects in this package
type Factory interface {
	FileStore() (Store, error)
	// GCalAuthHandler(fap *gcalauthParams, ftp *FitbitTokenParams) GCalAuthHandler
}