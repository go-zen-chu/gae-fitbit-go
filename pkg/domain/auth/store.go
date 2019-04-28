//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package auth

import "golang.org/x/oauth2"

// Store is an interface for storing oauth token
type Store interface {
	WriteFitbitToken(token *oauth2.Token) error
	FetchFitbitToken() (*oauth2.Token, error)
	WriteGCalToken(token *oauth2.Token) error
	FetchGCalToken() (*oauth2.Token, error)
}
