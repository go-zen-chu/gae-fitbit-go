package gcalauth

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
)

type oauthClient struct {
	config *oauth2.Config
}

func NewOAuthClient (config *oauth2.Config) dfba.OAuthClient {
	return &oauthClient{
		config: config,
	}
}

func (oc *oauthClient) GetAuthCodeURL() string {
	return oc.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (oc *oauthClient) Exchange(authCode string) (*oauth2.Token, error) {
	return oc.config.Exchange(context.Background(), authCode)
}