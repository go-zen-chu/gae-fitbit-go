//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package gcalauth

import "golang.org/x/oauth2"

// Store : responsible for storing data of fitbitauth
type Store interface {
	WriteGCalToken(token *oauth2.Token) error
	FetchGCalToken() (*oauth2.Token, error)
}
