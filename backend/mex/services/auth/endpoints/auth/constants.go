package auth

const (
	responseModeQuery    = "query"
	responseModeFragment = "fragment"
	signingKeyID         = "default"
	codeChallengeMethod  = "S256"
	responseTypeAuthCode = "code"

	grantAuthorizationCode = "authorization_code"
	grantClientCredentials = "client_credentials"
	grantRefreshToken      = "refresh_token"

	grantTypeAuthorizationCode = "0"
	grantTypeClientCredentials = "1"

	refreshTokenLength = 16
)
