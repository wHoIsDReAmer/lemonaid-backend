package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
)

func OAuthSetting() {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
