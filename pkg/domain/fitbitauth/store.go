package fitbitauth

// Store : responsible for storing data of fitbitauth
type Store interface {
	WriteAuthCode(authCode string) error
	FetchAuthCode() (string, error)
	WriteFitbitTokens(ft *FitbitTokens) error
	FetchFitbitTokens() (*FitbitTokens, error)
}
