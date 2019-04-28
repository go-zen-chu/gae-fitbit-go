package auth

import (
	"context"

	"golang.org/x/oauth2"
)

type OAuthClient interface {
	GetAuthCodeURL() string
	Exchange(authCode string) (*oauth2.Token, error)
}

type oauthClient struct {
	config *oauth2.Config
}

func NewOAuthClient(config *oauth2.Config) OAuthClient {
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
