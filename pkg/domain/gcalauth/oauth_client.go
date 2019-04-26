//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package gcalauth

import "golang.org/x/oauth2"

type OAuthClient interface {
	GetAuthCodeURL() string
	Exchange(authCode string) (*oauth2.Token, error)
}
