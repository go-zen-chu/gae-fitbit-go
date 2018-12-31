package fitbit2gcal

import (
	dfba "github.com/go-zen-chu/gae-fitbit-go/pkg/domain/fitbitauth"
	"golang.org/x/oauth2"
)

// Store : get stored fitbit token
type Store interface {
	FetchFitbitTokens() (*dfba.FitbitTokens, error)
	FetchGCalTokens() (*oauth2.Token, error)
}
