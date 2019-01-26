//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

import "golang.org/x/oauth2"

// Store : responsible for storing data of fitbitauth
type Store interface {
	WriteFitbitToken(token *oauth2.Token) error
	FetchFitbitToken() (*oauth2.Token, error)
}
