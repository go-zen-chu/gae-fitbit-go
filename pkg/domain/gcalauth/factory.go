//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package gcalauth

import "golang.org/x/oauth2"

// Factory : Creates objects in this package
type Factory interface {
	FileStore() (Store, error)
	GCalAuthHandler(oauthConfig *oauth2.Config) GCalAuthHandler
}
