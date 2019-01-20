//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

type FitbitHTTPClient interface {
	GetFitbitToken(fbtp *FitbitTokenParams, authCode string) (*FitbitTokens, error)
}
