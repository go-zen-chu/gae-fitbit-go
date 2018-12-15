package fitbitauth

// FitbitAuthConfig : [Using OAuth 2.0](https://dev.fitbit.com/build/reference/web-api/oauth2/#authorization-page)
type FitbitAuthConfig struct {
	ClientID     string
	Scope        string
	RedirectURI  string
	ResponseType string
	Expires      string
}
