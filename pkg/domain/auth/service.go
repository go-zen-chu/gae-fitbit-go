//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package auth

import (
	"github.com/pkg/errors"
)

type Service interface {
	HandleFitbitAuthCode(code string) error
	HandleGCalAuthCode(code string) error
}

type service struct {
	store             Store
	fitbitOAuthClient OAuthClient
	gcalOAuthClient   OAuthClient
}

func NewService(store Store, fitbitOAuthClient, gcalOAuthClient OAuthClient) Service {
	return &service{
		store:             store,
		fitbitOAuthClient: fitbitOAuthClient,
		gcalOAuthClient:   gcalOAuthClient,
	}
}

func (s *service) HandleFitbitAuthCode(code string) error {
	token, err := s.fitbitOAuthClient.Exchange(code)
	if err != nil {
		return errors.Wrap(err, "Error while getting token")
	}
	err = s.store.WriteFitbitToken(token)
	if err != nil {
		return errors.Wrap(err, "Error while storing token")
	}
	return nil
}

func (s *service) HandleGCalAuthCode(code string) error {
	token, err := s.gcalOAuthClient.Exchange(code)
	if err != nil {
		return errors.Wrap(err, "Error while getting token")
	}
	err = s.store.WriteGCalToken(token)
	if err != nil {
		return errors.Wrap(err, "Error while storing token")
	}
	return nil
}
