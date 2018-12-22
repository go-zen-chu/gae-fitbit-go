package fitbitauth

// FitbitAuthParams : [Using OAuth 2.0](https://dev.fitbit.com/build/reference/web-api/oauth2/#authorization-page)
type FitbitAuthParams struct {
	ClientID     string
	Scope        string
	RedirectURI  string
	ResponseType string
	Expires      string
}

// FitbitTokenParams : [Using OAuth 2.0](https://dev.fitbit.com/build/reference/web-api/oauth2/#authorization-page)
type FitbitTokenParams struct {
	ClientID    string
	GrantType   string
	RedirectURI string
}

// FitbitTokens : [Using OAuth 2.0](https://dev.fitbit.com/build/reference/web-api/oauth2/#access-token-request)
type FitbitTokens struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}
