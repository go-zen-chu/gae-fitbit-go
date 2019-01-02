//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

// Factory : Creates objects in this package
type Factory interface {
	FileStore() (Store, error)
	FitbitAuthHandler(fap *FitbitAuthParams, ftp *FitbitTokenParams) FitbitAuthHandler
}
