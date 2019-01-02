//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package fitbitauth

// Store : responsible for storing data of fitbitauth
type Store interface {
	WriteAuthCode(authCode string) error
	FetchAuthCode() (string, error)
	WriteFitbitTokens(ft *FitbitTokens) error
	FetchFitbitTokens() (*FitbitTokens, error)
}
