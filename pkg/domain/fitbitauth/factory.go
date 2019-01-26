//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

import (
	"golang.org/x/oauth2"
)

// Factory : Creates objects in this package
type Factory interface {
	FileStore() (Store, error)
	FitbitAuthHandler(config *oauth2.Config) FitbitAuthHandler
	OAuthClient(config *oauth2.Config) OAuthClient
}
